package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chenhg5/collection"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	controllers "api/app/http/controllers/api/v1"
	"api/pkg/book"
	"api/pkg/cache"
	"api/pkg/file"
	"api/pkg/jwt"
	"api/pkg/logger"
	"api/pkg/openai"
	"api/pkg/response"
)

// Manager 所有 websocket 信息
type Manager struct {
	Group                   map[string]map[string]*Client
	groupCount, clientCount uint
	Lock                    sync.Mutex
	Register, UnRegister    chan *Client
	Message                 chan *MessageData
	GroupMessage            chan *GroupMessageData
	BroadCastMessage        chan *BroadCastMessageData
}

// Client 单个 websocket 信息
type Client struct {
	Id, Group  string
	Context    context.Context
	CancelFunc context.CancelFunc
	Socket     *websocket.Conn
	Message    chan []byte
}

// MessageData 单个发送数据信息
type MessageData struct {
	Id, Group string
	Context   context.Context
	Message   []byte
}

// GroupMessageData 组广播数据信息
type GroupMessageData struct {
	Group   string
	Message []byte
}

// BroadCastMessageData 广播发送数据信息
type BroadCastMessageData struct {
	Message []byte
}

//RequestParameters 请求参数
type RequestParameters struct {
	Group      string            `json:"group"`
	Parameters map[string]string `json:"parameters"`
	Chat       Chat              `json:"chat"`
}

// Volume 卷
type Chat struct {
	Content string
	Time    int
	Type    string
	From    struct {
		AvatarUrl string `json:"avatarUrl"`
		Id        string `json:"id"`
		Name      string `json:"name"`
	}
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read(cxt context.Context) {
	defer func(cxt context.Context) {
		WebsocketManager.UnRegister <- c
		logger.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			logger.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}(cxt)

	for {
		if cxt.Err() != nil {
			break
		}
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}
		content := &RequestParameters{}
		if err = json.Unmarshal(message, &content); err == nil {
			if content.Group == "book" {
				if content.Parameters["url"] != "" {
					book.Download(cxt, content.Parameters["url"], c.Id, c.Group, SendOne)
				}
				if content.Parameters["log"] != "" {
					book.DownloadLog(cxt, content.Parameters["log"], c.Id, c.Group, SendOne)
				}
			}
			if content.Group == "chat" {
				chat := cache.Get("chat")
				b, _ := json.Marshal(&chat)
				var m []Chat
				_ = json.Unmarshal(b, &m)
				cache.Set("chat", append(m, content.Chat), time.Hour*24*30*12)
				list, _ := json.Marshal(content.Chat)
				SendGroup(list, content.Group, "message")
				if strings.Contains(content.Chat.Content, "@ai") {
					msg := strings.Replace(content.Chat.Content, "@ai", "", 1)
					gpt := openai.NewChatGptTool("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")
					message := []openai.Gpt3Dot5Message{
						{
							Role:    "user",
							Content: msg,
						},
					}
					res, err := gpt.ChatGPT3Dot5Turbo(message)
					if err == nil {
						var ai = Chat{
							Content: res,
							Time:    0,
							Type:    "text",
							From: struct {
								AvatarUrl string `json:"avatarUrl"`
								Id        string `json:"id"`
								Name      string `json:"name"`
							}(struct {
								AvatarUrl string
								Id        string
								Name      string
							}{AvatarUrl: "https://xsgames.co/randomusers/assets/avatars/pixel/0.jpg", Id: "ai", Name: "机器人"}),
						}
						chat := cache.Get("chat")
						b, _ := json.Marshal(&chat)
						var m []Chat
						_ = json.Unmarshal(b, &m)
						cache.Set("chat", append(m, ai), time.Hour*24*30*12)
						list, _ := json.Marshal(ai)
						SendGroup(list, content.Group, "message")
					}
				}
			}
		}
		logger.Printf("client [%s] receive message: %s", c.Id, string(message))
		c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write(cxt context.Context) {
	defer func(cxt context.Context) {
		logger.Printf("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			logger.Printf("client [%s] disconnect err: %s", c.Id, err)
		}
	}(cxt)

	for {
		if cxt.Err() != nil {
			break
		}
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			logger.Printf("client [%s] write message: %s", c.Id, string(message))
			if string(message) == "monitor" {
				message = new(controllers.ServerController).GetServerInfo()
			}
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Printf("client [%s] writemessage err: %s", c.Id, err)
			}
		case _ = <-c.Context.Done():
			break
		}
	}
}

// Start 启动 websocket 管理器
func (manager *Manager) Start() {
	logger.Printf("websocket manage start")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			logger.Printf("client [%s] connect", client.Id)
			logger.Printf("register client [%s] to group [%s]", client.Id, client.Group)

			manager.Lock.Lock()
			if manager.Group[client.Group] == nil {
				manager.Group[client.Group] = make(map[string]*Client)
				manager.groupCount += 1
			}
			manager.Group[client.Group][client.Id] = client
			manager.clientCount += 1
			manager.Lock.Unlock()

		// 注销
		case client := <-manager.UnRegister:
			logger.Printf("unregister client [%s] from group [%s]", client.Id, client.Group)
			manager.Lock.Lock()
			if mGroup, ok := manager.Group[client.Group]; ok {
				if mClient, ok := mGroup[client.Id]; ok {
					close(mClient.Message)
					delete(mGroup, client.Id)
					manager.clientCount -= 1
					if len(mGroup) == 0 {
						//logger.Printf("delete empty group [%s]", client.Group)
						delete(manager.Group, client.Group)
						manager.groupCount -= 1
					}
					mClient.CancelFunc()
				}
			}
			manager.Lock.Unlock()

			// 发送广播数据到某个组的 channel 变量 Send 中
			//case data := <-manager.boardCast:
			//	if groupMap, ok := manager.wsGroup[data.GroupId]; ok {
			//		for _, conn := range groupMap {
			//			conn.Send <- data.Data
			//		}
			//	}
		}
	}
}

// SendService 处理单个 client 发送数据
func (manager *Manager) SendService() {
	for {
		select {
		case data := <-manager.Message:
			if groupMap, ok := manager.Group[data.Group]; ok {
				if conn, ok := groupMap[data.Id]; ok {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendGroupService 处理 group 广播数据
func (manager *Manager) SendGroupService() {
	for {
		select {
		// 发送广播数据到某个组的 channel 变量 Send 中
		case data := <-manager.GroupMessage:
			if groupMap, ok := manager.Group[data.Group]; ok {
				for _, conn := range groupMap {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// SendAllService 处理广播数据
func (manager *Manager) SendAllService() {
	for {
		select {
		case data := <-manager.BroadCastMessage:
			for _, v := range manager.Group {
				for _, conn := range v {
					conn.Message <- data.Message
				}
			}
		}
	}
}

// Send 向指定的 client 发送数据
func (manager *Manager) Send(cxt context.Context, id string, group string, message []byte) {
	data := &MessageData{
		Id:      id,
		Context: cxt,
		Group:   group,
		Message: message,
	}
	manager.Message <- data
}

// SendGroup 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, message []byte) {
	data := &GroupMessageData{
		Group:   group,
		Message: message,
	}
	manager.GroupMessage <- data
}

// SendAll 广播
func (manager *Manager) SendAll(message []byte) {
	data := &BroadCastMessageData{
		Message: message,
	}
	manager.BroadCastMessage <- data
}

// RegisterClient 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnRegisterClient 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// LenGroup 当前组个数
func (manager *Manager) LenGroup() uint {
	return manager.groupCount
}

// LenClient 当前连接个数
func (manager *Manager) LenClient() uint {
	return manager.clientCount
}

// Info 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	managerInfo["chanMessageLen"] = len(manager.Message)
	managerInfo["chanGroupMessageLen"] = len(manager.GroupMessage)
	managerInfo["chanBroadCastMessageLen"] = len(manager.BroadCastMessage)
	return managerInfo
}

// WebsocketManager 初始化 wsManager 管理器
var WebsocketManager = Manager{
	Group:            make(map[string]map[string]*Client),
	Register:         make(chan *Client, 128),
	UnRegister:       make(chan *Client, 128),
	GroupMessage:     make(chan *GroupMessageData, 128),
	Message:          make(chan *MessageData, 128),
	BroadCastMessage: make(chan *BroadCastMessageData, 128),
	groupCount:       0,
	clientCount:      0,
}

var typeWhitelist = []string{"book", "chat"}

// WsClient gin 处理 websocket handler
func (manager *Manager) WsClient(c *gin.Context) {

	ctx, cancel := context.WithCancel(context.Background())

	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	channel := c.Param("channel")
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		cancel()
		logger.Printf("websocket connect error: %s", channel)
		return
	}

	token := c.Query("token")
	var claims *jwt.JWTCustomClaims
	if !collection.Collect(typeWhitelist).Contains(channel) {
		claims, err = jwt.NewJWT().WsParserToken(token)
		// JWT 解析失败，有错误发生
		if err != nil {
			cancel()
			logger.Printf("tokenerror: %s", token)
			_ = conn.Close()
			return
		}
	}
	userId := token
	if claims != nil {
		userId = claims.UserID
	}
	client := &Client{
		Id:         userId,
		Group:      channel,
		Context:    ctx,
		CancelFunc: cancel,
		Socket:     conn,
		Message:    make(chan []byte, 1024),
	}

	manager.RegisterClient(client)
	go client.Read(ctx)
	go client.Write(ctx)
	if channel == "cron" {
		file.FileMonitoring(ctx, "storage/cron/"+time.Now().Format("2006-01-02.log"), userId, channel, SendOne)
	}
	if channel == "chat" {
		time.Sleep(200 * time.Millisecond)

		info, _ := json.Marshal(WebsocketManager.Info())
		SendGroup(info, channel, "info")
		chat := cache.Get("chat")
		list, _ := json.Marshal(chat)
		SendOne(ctx, userId, channel, list)
	}
}

func (manager *Manager) UnWsClient(c *gin.Context) {
	token := c.Query("token")
	group := c.Param("channel")
	var claims *jwt.JWTCustomClaims
	if !collection.Collect(typeWhitelist).Contains(group) {
		err := error(nil)
		claims, err = jwt.NewJWT().WsParserToken(c.Query("token"))
		if err != nil {
			response.NormalVerificationError(c, "退出失败")
			return
		}
	}
	id := token
	if claims != nil {
		id = claims.UserID
	}
	Logout(id, group)
	c.Set("result", "ws close success")
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": "ws close success",
		"msg":  "success",
	})
}

func SendGroup(msg []byte, group string, t string) {
	WebsocketManager.SendGroup(group, []byte("{\"code\":200,\"type\":\""+t+"\",\"data\":"+string(msg)+"}"))
	logger.Dump(WebsocketManager.Info())
}

func SendAll(msg []byte) {
	WebsocketManager.SendAll([]byte("{\"code\":200,\"data\":" + string(msg) + "}"))
	logger.Dump(WebsocketManager.Info())
}

func SendOne(ctx context.Context, id string, group string, msg []byte) {
	WebsocketManager.Send(ctx, id, group, []byte("{\"code\":200,\"data\":"+string(msg)+"}"))
	logger.Dump(WebsocketManager.Info())
}

func Logout(id string, group string) {
	WebsocketManager.UnRegisterClient(&Client{Id: id, Group: group})
	logger.Dump(WebsocketManager.Info())
}
