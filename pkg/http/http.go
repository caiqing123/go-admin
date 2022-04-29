package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Get 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)
	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), nil
}

// Post 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data io.Reader, contentType string) ([]byte, error) {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	//jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, data)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}
