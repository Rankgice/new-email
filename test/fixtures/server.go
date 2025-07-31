package fixtures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"new-email/internal/router"
	"strings"

	"github.com/gin-gonic/gin"
)

// TestServer 测试服务器结构
type TestServer struct {
	Router *gin.Engine
	TestDB *TestDB
}

// NewTestServer 创建测试服务器
func NewTestServer() *TestServer {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试数据库
	testDB := NewTestDB()

	// 创建Gin引擎
	r := gin.New()

	// 设置路由
	router.SetupRouter(r, testDB.SvcCtx)

	return &TestServer{
		Router: r,
		TestDB: testDB,
	}
}

// Close 关闭测试服务器
func (ts *TestServer) Close() {
	ts.TestDB.Close()
}

// Clean 清理测试数据
func (ts *TestServer) Clean() {
	ts.TestDB.Clean()
}

// SeedData 插入测试数据
func (ts *TestServer) SeedData() {
	ts.TestDB.SeedTestData()
}

// Request 发送HTTP请求的辅助方法
func (ts *TestServer) Request(method, path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	var reqBody io.Reader

	// 处理请求体
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = strings.NewReader(v)
		case []byte:
			reqBody = bytes.NewReader(v)
		default:
			jsonData, _ := json.Marshal(v)
			reqBody = bytes.NewReader(jsonData)
		}
	}

	// 创建请求
	req, _ := http.NewRequest(method, path, reqBody)

	// 设置默认Content-Type
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置自定义头部
	for _, header := range headers {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	ts.Router.ServeHTTP(w, req)

	return w
}

// GET 发送GET请求
func (ts *TestServer) GET(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	return ts.Request("GET", path, nil, headers...)
}

// POST 发送POST请求
func (ts *TestServer) POST(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	return ts.Request("POST", path, body, headers...)
}

// PUT 发送PUT请求
func (ts *TestServer) PUT(path string, body interface{}, headers ...map[string]string) *httptest.ResponseRecorder {
	return ts.Request("PUT", path, body, headers...)
}

// DELETE 发送DELETE请求
func (ts *TestServer) DELETE(path string, headers ...map[string]string) *httptest.ResponseRecorder {
	return ts.Request("DELETE", path, nil, headers...)
}

// ParseResponse 解析响应JSON
func (ts *TestServer) ParseResponse(w *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(w.Body.Bytes(), v)
}

// AssertStatus 断言HTTP状态码
func (ts *TestServer) AssertStatus(w *httptest.ResponseRecorder, expectedStatus int) bool {
	if w.Code != expectedStatus {
		fmt.Printf("期望状态码 %d，实际状态码 %d\n", expectedStatus, w.Code)
		fmt.Printf("响应内容: %s\n", w.Body.String())
		return false
	}
	return true
}

// AssertJSON 断言响应是有效的JSON
func (ts *TestServer) AssertJSON(w *httptest.ResponseRecorder) bool {
	var js json.RawMessage
	return json.Unmarshal(w.Body.Bytes(), &js) == nil
}

// AssertContains 断言响应包含指定内容
func (ts *TestServer) AssertContains(w *httptest.ResponseRecorder, content string) bool {
	return strings.Contains(w.Body.String(), content)
}

// GetResponseData 获取响应中的data字段
func (ts *TestServer) GetResponseData(w *httptest.ResponseRecorder) (interface{}, error) {
	var response struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	err := ts.ParseResponse(w, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// IsSuccess 检查响应是否成功
func (ts *TestServer) IsSuccess(w *httptest.ResponseRecorder) bool {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err := ts.ParseResponse(w, &response)
	if err != nil {
		return false
	}

	return response.Code == 0
}
