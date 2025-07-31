package api

import (
	"fmt"
	"net/http"
	"new-email/test/fixtures"
	"testing"
)

// TestJWTDebug 调试JWT验证问题
func TestJWTDebug(t *testing.T) {
	server := fixtures.NewTestServer()
	defer server.Close()

	// 先注册一个用户
	registerData := map[string]interface{}{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "password123",
		"nickname": "测试用户",
	}

	w1 := server.POST("/api/public/user/register", registerData)
	if !server.AssertStatus(w1, http.StatusOK) {
		t.Fatalf("注册失败")
	}

	// 登录获取token
	loginData := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	if !server.AssertStatus(w2, http.StatusOK) {
		t.Fatalf("登录失败")
	}

	var loginResponse struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
			User  struct {
				Id       int64  `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"data"`
	}

	err := server.ParseResponse(w2, &loginResponse)
	if err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}

	token := loginResponse.Data.Token
	userId := loginResponse.Data.User.Id

	t.Logf("登录成功: UserId=%d, Token=%s", userId, token[:50]+"...")

	// 测试需要认证的接口
	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// 测试用户信息接口
	w3 := server.GET("/api/user/profile", authHeader)
	t.Logf("用户信息接口响应状态: %d", w3.Code)
	t.Logf("用户信息接口响应内容: %s", w3.Body.String())

	if w3.Code != 200 {
		t.Fatalf("用户信息接口调用失败")
	}

	// 测试邮箱列表接口
	w4 := server.GET("/api/user/mailboxes", authHeader)
	t.Logf("邮箱列表接口响应状态: %d", w4.Code)
	t.Logf("邮箱列表接口响应内容: %s", w4.Body.String())

	if w4.Code != 200 {
		t.Fatalf("邮箱列表接口调用失败")
	}

	t.Logf("✅ JWT验证测试通过")
}
