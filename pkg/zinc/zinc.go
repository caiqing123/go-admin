package zinc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ZincClient struct {
	*ZincClientConfig
}

type ZincClientConfig struct {
	ZincHost     string
	ZincUser     string
	ZincPassword string
}

type ZincIndex struct {
	Name        string             `json:"name"`
	StorageType string             `json:"storage_type"`
	Mappings    *ZincIndexMappings `json:"mappings"`
}

type ZincIndexMappings struct {
	Properties *ZincIndexProperty `json:"properties"`
}

type ZincIndexProperty map[string]*ZincIndexPropertyT

type ZincIndexPropertyT struct {
	Type           string `json:"type"`
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`
	Analyzer       string `json:"analyzer"`
	SearchAnalyzer string `json:"search_analyzer"`
	Format         string `json:"format"`
}

type QueryResultT struct {
	Took     int          `json:"took"`
	TimedOut bool         `json:"timed_out"`
	Hits     *HitsResultT `json:"hits"`
}
type HitsResultT struct {
	Total    *HitsResultTotalT `json:"total"`
	MaxScore float64           `json:"max_score"`
	Hits     []*HitItem        `json:"hits"`
}
type HitsResultTotalT struct {
	Value int64 `json:"value"`
}

type HitItem struct {
	Index     string      `json:"_index"`
	Type      string      `json:"_type"`
	ID        string      `json:"_id"`
	Score     float64     `json:"_score"`
	Timestamp time.Time   `json:"@timestamp"`
	Source    interface{} `json:"_source"`
}

type IndexResult []*IndexResultS

type IndexResultS struct {
	Name        string      `json:"name"`
	StorageType string      `json:"storage_type"`
	Settings    interface{} `json:"settings"`
	Mappings    interface{} `json:"mappings"`
	CreateAt    string      `json:"create_at"`
	UpdateAt    string      `json:"update_at"`
	DocsCount   int         `json:"docs_count"`
	StorageSize int         `json:"storage_size"`
}

// CreateIndex 创建索引
func (c *ZincClient) CreateIndex(name string, p *ZincIndexProperty) bool {
	data := &ZincIndex{
		Name:        name,
		StorageType: "disk",
		Mappings: &ZincIndexMappings{
			Properties: p,
		},
	}
	prtBodyBytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return false
	}
	_, err = c.request("PUT", "/api/index", strings.NewReader(string(prtBodyBytes)))
	if err != nil {
		return false
	}
	return true
}

// ExistIndex 检查索引是否存在
func (c *ZincClient) ExistIndex(name string) bool {
	resp, err := c.request("GET", "/api/index", nil)
	if err != nil {
		return false
	}
	retData := &IndexResult{}
	err = json.Unmarshal(resp, retData)

	if err != nil {
		return false
	}
	for _, v := range *retData {
		if v.Name == name {
			return true
		}
	}
	return false
}

// PutDoc 新增/更新文档
func (c *ZincClient) PutDoc(name string, id int64, doc interface{}) (bool, error) {
	prtBodyBytes, err := json.MarshalIndent(doc, "", "   ")
	if err != nil {
		return false, err
	}
	_, err = c.request("PUT", fmt.Sprintf("/api/%s/_doc/%d", name, id), strings.NewReader(string(prtBodyBytes)))
	if err != nil {
		return false, err
	}

	return true, nil
}

// BulkPutLogDoc 批量新增文档 日志使用
func (c *ZincClient) BulkPutLogDoc(docs []map[string]interface{}) (bool, error) {
	dataStr := ""
	for _, doc := range docs {
		str, err := json.Marshal(doc)
		if err == nil {
			dataStr = dataStr + string(str) + "\n"
		}
	}
	_, err := c.request("POST", "/es/_bulk", strings.NewReader(dataStr))
	if err != nil {
		return false, err
	}

	return true, nil
}

// BulkPushDoc 批量新增文档
func (c *ZincClient) BulkPushDoc(docs []map[string]interface{}) (bool, error) {
	dataStr := ""
	for _, doc := range docs {
		str, err := json.Marshal(doc)
		if err == nil {
			dataStr = dataStr + string(str) + "\n"
		}
	}
	_, err := c.request("POST", "/api/_bulk", strings.NewReader(dataStr))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *ZincClient) EsQuery(indexName string, q interface{}) (*QueryResultT, error) {
	prtBodyBytes, err := json.MarshalIndent(q, "", "   ")
	if err != nil {
		return nil, err
	}
	resp1, err := c.request("POST", fmt.Sprintf("/es/%s/_search", indexName), strings.NewReader(string(prtBodyBytes)))
	if err != nil {
		return nil, err
	}
	result := &QueryResultT{}
	err = json.Unmarshal(resp1, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ZincClient) ApiQuery(indexName string, q interface{}) (*QueryResultT, error) {
	prtBodyBytes, err := json.MarshalIndent(q, "", "   ")
	if err != nil {
		return nil, err
	}
	resp1, err := c.request("POST", fmt.Sprintf("/api/%s/_search", indexName), strings.NewReader(string(prtBodyBytes)))
	if err != nil {
		return nil, err
	}
	result := &QueryResultT{}
	err = json.Unmarshal(resp1, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ZincClient) DelDoc(indexName, id string) error {
	_, err := c.request("DELETE", fmt.Sprintf("/api/%s/_doc/%s", indexName, id), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ZincClient) request(method, url string, body io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, c.ZincHost+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.ZincUser+":"+c.ZincPassword)))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	resp1, _ := ioutil.ReadAll(resp.Body)
	return resp1, nil
}
