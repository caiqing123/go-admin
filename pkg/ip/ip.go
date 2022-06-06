package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetLocation 获取外网ip地址
func GetLocation(ip string) string {
	if ip == "127.0.0.1" || ip == "localhost" {
		return "内部IP"
	}
	key := "d3af6d5474dea1410b613a1f4e223e8b"
	url := "https://restapi.amap.com/v3/ip?ip=" + ip + "&key=" + key
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("restapi.amap.com failed:", err)
		return "未知位置"
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	s, err := ioutil.ReadAll(resp.Body)

	m := make(map[string]string)

	err = json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
	}
	if m["status"] == "0" {
		return m["info"]
	}
	return m["country"] + "-" + m["province"] + "-" + m["city"] + "-" + m["district"] + "-" + m["isp"]
}

// GetLocationHost 获取局域网ip地址
func GetLocationHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}

func GetClientIP(c *gin.Context) string {
	ClientIP := c.ClientIP()
	//fmt.Println("ClientIP:", ClientIP)
	RemoteIP, _ := c.RemoteIP()
	//fmt.Println("RemoteIP:", RemoteIP)
	ip := c.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	if RemoteIP.String() != "127.0.0.1" {
		ip = RemoteIP.String()
	}
	if ClientIP != "127.0.0.1" {
		ip = ClientIP
	}
	return ip
}
