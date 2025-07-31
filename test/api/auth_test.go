package api

import (
	"net/http"
	"new-email/test/fixtures"
	"testing"
)

// TestUserRegister 测试用户注册
func TestUserRegister(t *testing.T) {
	// 创建测试服务器
	server := fixtures.NewTestServer()
	defer server.Close()

	// 测试数据
	registerData := map[string]interface{}{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "password123",
		"nickname": "测试用户",
	}

	// 发送注册请求
	w := server.POST("/api/public/user/register", registerData)

	// 验证响应状态码
	if !server.AssertStatus(w, http.StatusOK) {
		t.Fatalf("注册请求失败")
	}

	// 验证响应格式
	if !server.AssertJSON(w) {
		t.Fatalf("响应不是有效的JSON格式")
	}

	// 解析响应
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id       int64  `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}

	err := server.ParseResponse(w, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 验证响应内容
	if response.Code != 0 {
		t.Fatalf("注册失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	if response.Data.Username != "testuser" {
		t.Fatalf("用户名不匹配: 期望 testuser, 实际 %s", response.Data.Username)
	}

	if response.Data.Email != "testuser@example.com" {
		t.Fatalf("邮箱不匹配: 期望 testuser@example.com, 实际 %s", response.Data.Email)
	}

	t.Logf("✅ 用户注册测试通过: ID=%d, Username=%s", response.Data.Id, response.Data.Username)
}

// TestUserRegisterDuplicate 测试重复注册
func TestUserRegisterDuplicate(t *testing.T) {
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
		t.Fatalf("第一次注册失败")
	}

	// 再次注册相同用户名
	w2 := server.POST("/api/public/user/register", registerData)
	if !server.AssertStatus(w2, http.StatusOK) {
		t.Fatalf("第二次注册请求失败")
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err := server.ParseResponse(w2, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 应该返回错误
	if response.Code == 0 {
		t.Fatalf("重复注册应该失败，但返回成功")
	}

	if !server.AssertContains(w2, "用户名已存在") {
		t.Fatalf("错误信息不正确，应该包含'用户名已存在'")
	}

	t.Logf("✅ 重复注册测试通过: %s", response.Msg)
}

// TestUserLogin 测试用户登录
func TestUserLogin(t *testing.T) {
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

	// 测试登录
	loginData := map[string]interface{}{
		"username": "testuser",
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	if !server.AssertStatus(w2, http.StatusOK) {
		t.Fatalf("登录请求失败")
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token string `json:"token"`
			User  struct {
				Id       int64  `json:"id"`
				Username string `json:"username"`
				Email    string `json:"email"`
				Nickname string `json:"nickname"`
				Status   int    `json:"status"`
			} `json:"user"`
		} `json:"data"`
	}

	err := server.ParseResponse(w2, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("登录失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	if response.Data.Token == "" {
		t.Fatalf("登录成功但未返回token")
	}

	if response.Data.User.Username != "testuser" {
		t.Fatalf("用户信息不匹配")
	}

	t.Logf("✅ 用户登录测试通过: Token=%s...", response.Data.Token[:20])
}

// TestUserLoginWithEmail 测试使用邮箱登录
func TestUserLoginWithEmail(t *testing.T) {
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

	// 使用邮箱登录
	loginData := map[string]interface{}{
		"username": "testuser@example.com", // 使用邮箱作为用户名
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	if !server.AssertStatus(w2, http.StatusOK) {
		t.Fatalf("邮箱登录请求失败")
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	err := server.ParseResponse(w2, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("邮箱登录失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	if response.Data.Token == "" {
		t.Fatalf("邮箱登录成功但未返回token")
	}

	t.Logf("✅ 邮箱登录测试通过")
}

// TestUserLoginWrongPassword 测试错误密码登录
func TestUserLoginWrongPassword(t *testing.T) {
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

	// 使用错误密码登录
	loginData := map[string]interface{}{
		"username": "testuser",
		"password": "wrongpassword",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	if !server.AssertStatus(w2, http.StatusOK) {
		t.Fatalf("错误密码登录请求失败")
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err := server.ParseResponse(w2, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	// 应该返回错误
	if response.Code == 0 {
		t.Fatalf("错误密码登录应该失败，但返回成功")
	}

	t.Logf("✅ 错误密码登录测试通过: %s", response.Msg)
}
