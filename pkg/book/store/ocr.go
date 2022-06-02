package store

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const (
	requestMethod = "POST"
	service       = "ocr"
	host          = "ocr.tencentcloudapi.com"
	region        = "ap-beijing"
	algorithm     = "TC3-HMAC-SHA256"
	contentType   = "application/json"
	tc3           = "tc3_request"
	signedHeaders = "content-type;host"
	RequestUrl    = "https://ocr.tencentcloudapi.com"
)

type OcrMethodData struct {
	Name    string
	Version string
}

var GeneralMethod = OcrMethodData{"GeneralEfficientOCR", "2018-11-19"}

type GeneralData struct {
	Response struct {
		Language       string `json:"language"`
		RequestId      string `json:"requestId"`
		TextDetections []struct {
			DetectedText string `json:"detectedText"` // 识别出的文本行内容
			Confidence   int    `json:"confidence"`   // 置信度 0 ~100
			Polygon      []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"polygon"` // 文本行坐标，以四个顶点坐标表示 注意：此字段可能返回 null，表示取不到有效值。
			AdvancedInfo string `json:"advancedInfo"` //此字段为扩展字段。 GeneralBasicOcr接口返回段落信息Parag ，包含ParagNo。
		} `json:"textDetections"`
		Error struct {
			Code    string
			Message string
		}
	}
}

func (generalData GeneralData) Error() bool {
	return generalData.Response.Error.Code != ""
}

func Request(ocr OcrMethodData, img []byte) (*http.Response, error) {
	params := make(map[string]string, 2)
	if ocr.Name == "IDCardOCR" {
		params["CardSide"] = "FRONT" //FRONT为身份证有照片的一面（正面） BACK为身份证有国徽的一面（反面）
		//params["Config"] = ""   //可选字段，根据需要选择是否请求对应字段。目前包含的字段为： CropIdCard-身份证照片裁剪， CropPortrait-人像照片裁剪， CopyWarn-复印件告警， ReshootWarn-翻拍告警。
	}

	data := base64.StdEncoding.EncodeToString(img)
	params["ImageBase64"] = data
	client := &http.Client{}
	js, _ := json.Marshal(params)
	//fmt.Printf("params: %s\n\n", string(js))
	req, _ := http.NewRequest(http.MethodPost, RequestUrl, bytes.NewReader(js))
	req.Header = getHeader(js, ocr.Name, ocr.Version)
	return client.Do(req)
}

// 获取请求的 Header
func getHeader(payload []byte, action, version string) http.Header {
	header := http.Header{}
	tm := time.Now()
	timeSpan := strconv.FormatInt(tm.Unix(), 10)
	header.Add("Authorization", getAuth(tm, payload))
	header.Add("Host", host)
	header.Add("Content-Type", contentType)
	header.Add("X-TC-Action", action)
	header.Add("X-TC-Timestamp", timeSpan)
	header.Add("X-TC-Version", version)
	header.Add("X-TC-Region", region)
	return header
}

// 1. 拼接规范请求串
func getSpliceData(payload []byte) string {
	h := sha256.New()
	h.Write(payload)
	hs := h.Sum(nil)

	httpRequestMethod := requestMethod
	uri := "/"
	query := ""
	headers := "content-type:" + contentType + "\nhost:" + host + "\n"
	hashedRequestPayload := hex.EncodeToString(hs)

	request := httpRequestMethod + "\n" +
		uri + "\n" +
		query + "\n" +
		headers + "\n" +
		signedHeaders + "\n" +
		hashedRequestPayload

	//fmt.Printf("getSpliceData result: %s\n\n", request)
	return request
}

// 2. 拼接待签名字符串
func getStringToSign(splice string, tm time.Time) string {
	h := sha256.New()
	h.Write([]byte(splice))
	hs1 := h.Sum(nil)

	hashedRequest := hex.EncodeToString(hs1)
	scope := tm.Format("2006-01-02") + "/" + service + "/" + tc3
	timeSpan := strconv.FormatInt(tm.Unix(), 10)

	stringToSign := algorithm + "\n" + timeSpan + "\n" + scope + "\n" + hashedRequest
	//fmt.Printf("getStringToSign result: %s\n\n", stringToSign)
	return stringToSign
}

// 	3. 计算签名
func getSignature(stringToSign string, tm time.Time) string {
	secretDate := hmac256Secret([]byte("TC3wd7MM7csvkMjrkOnyub4JI1VjjspwpTt"), tm.Format("2006-01-02"))
	secretService := hmac256Secret(secretDate, service)
	secretSigning := hmac256Secret(secretService, tc3)
	signature := hex.EncodeToString(hmac256Secret(secretSigning, stringToSign))
	//fmt.Printf("getSignature result: %s\n\n", signature)
	return signature
}

// 4. 拼接 Authorization
func getAuth(tm time.Time, payload []byte) string {
	scope := tm.Format("2006-01-02") + "/" + service + "/" + tc3
	auth := algorithm + " " + "Credential=AKIDMWS8ZkrKirkOnxFWv62kDRQAf061h6GI/" + scope + ", " +
		"SignedHeaders=" + signedHeaders + ", " + "Signature=" +
		getSignature(getStringToSign(getSpliceData(payload), tm), tm)
	//fmt.Printf("auth result: %s\n\n", auth)
	return auth
}

// 哈希256摘要
func hmac256Secret(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	sha := mac.Sum(nil)
	return sha
}
