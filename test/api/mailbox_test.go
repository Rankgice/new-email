package api

import (
	"fmt"
	"net/http"
	"new-email/test/fixtures"
	"testing"
)

// TestMailboxCreate 测试创建邮箱
func TestMailboxCreate(t *testing.T) {
	server := fixtures.NewTestServer()
	defer server.Close()

	// 先插入域名数据
	server.SeedData()

	// 注册一个新用户（避免与SeedData中的用户冲突）
	registerData := map[string]interface{}{
		"username": "newtestuser",
		"email":    "newtestuser@example.com",
		"password": "password123",
		"nickname": "新测试用户",
	}

	w1 := server.POST("/api/public/user/register", registerData)
	if !server.AssertStatus(w1, http.StatusOK) {
		t.Fatalf("注册失败")
	}

	// 登录获取token
	loginData := map[string]interface{}{
		"username": "newtestuser",
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
		} `json:"data"`
	}

	err := server.ParseResponse(w2, &loginResponse)
	if err != nil {
		t.Fatalf("解析登录响应失败: %v", err)
	}

	token := loginResponse.Data.Token
	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	// 创建邮箱（使用测试数据中的域名ID=1）
	mailboxData := map[string]interface{}{
		"domainId": 1,                  // 使用测试域名
		"email":    "newtest@test.com", // 使用不同的邮箱地址
		"password": "mailbox_password",
		"status":   1,
	}

	w3 := server.POST("/api/user/mailboxes", mailboxData, authHeader)
	if !server.AssertStatus(w3, http.StatusOK) {
		t.Fatalf("创建邮箱失败: %s", w3.Body.String())
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id    int64  `json:"id"`
			Email string `json:"email"`
		} `json:"data"`
	}

	err = server.ParseResponse(w3, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("创建邮箱失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	t.Logf("创建邮箱响应: %s", w3.Body.String())

	if response.Data.Email != "newtest@test.com" {
		t.Fatalf("邮箱地址不匹配: 期望 newtest@test.com, 实际 %s", response.Data.Email)
	}

	t.Logf("✅ 创建邮箱测试通过: ID=%d, Email=%s", response.Data.Id, response.Data.Email)
}

// TestMailboxList 测试邮箱列表
func TestMailboxList(t *testing.T) {
	server := fixtures.NewTestServer()
	defer server.Close()

	// 插入测试数据
	server.SeedData()

	// 先注册用户并登录
	registerData := map[string]interface{}{
		"username": "listuser",
		"email":    "listuser@example.com",
		"password": "password123",
		"nickname": "列表测试用户",
	}

	server.POST("/api/public/user/register", registerData)

	loginData := map[string]interface{}{
		"username": "listuser",
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	var loginResponse struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	server.ParseResponse(w2, &loginResponse)

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", loginResponse.Data.Token),
	}

	// 创建几个邮箱
	mailboxes := []string{
		"list1@test.com",
		"list2@test.com",
		"list3@test.com",
	}

	for _, email := range mailboxes {
		mailboxData := map[string]interface{}{
			"domainId": 1, // 使用测试域名
			"email":    email,
			"password": "mailbox_password",
			"status":   1,
		}
		server.POST("/api/user/mailboxes", mailboxData, authHeader)
	}

	// 获取邮箱列表
	w3 := server.GET("/api/user/mailboxes", authHeader)
	if !server.AssertStatus(w3, http.StatusOK) {
		t.Fatalf("获取邮箱列表失败: %s", w3.Body.String())
	}

	var response struct {
		Code int `json:"code"`
		Data struct {
			List     []map[string]interface{} `json:"list"`
			Total    int64                    `json:"total"`
			Page     int                      `json:"page"`
			PageSize int                      `json:"pageSize"`
		} `json:"data"`
	}

	err := server.ParseResponse(w3, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("获取邮箱列表失败: code=%d", response.Code)
	}

	if len(response.Data.List) != 3 {
		t.Fatalf("邮箱数量不正确: 期望 3, 实际 %d", len(response.Data.List))
	}

	if response.Data.Total != 3 {
		t.Fatalf("总数不正确: 期望 3, 实际 %d", response.Data.Total)
	}

	t.Logf("✅ 邮箱列表测试通过: 共 %d 个邮箱", len(response.Data.List))
}

// TestMailboxUpdate 测试更新邮箱
func TestMailboxUpdate(t *testing.T) {
	server := fixtures.NewTestServer()
	defer server.Close()

	// 插入测试数据
	server.SeedData()

	// 先注册用户并登录
	registerData := map[string]interface{}{
		"username": "updateuser",
		"email":    "updateuser@example.com",
		"password": "password123",
		"nickname": "更新测试用户",
	}

	server.POST("/api/public/user/register", registerData)

	loginData := map[string]interface{}{
		"username": "updateuser",
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	var loginResponse struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	server.ParseResponse(w2, &loginResponse)

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", loginResponse.Data.Token),
	}

	// 创建邮箱
	mailboxData := map[string]interface{}{
		"domainId": 1, // 使用测试域名
		"email":    "update@test.com",
		"password": "mailbox_password",
		"status":   1,
	}

	w3 := server.POST("/api/user/mailboxes", mailboxData, authHeader)
	var createResponse struct {
		Code int `json:"code"`
		Data struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}
	server.ParseResponse(w3, &createResponse)

	mailboxId := createResponse.Data.Id

	// 更新邮箱
	updateData := map[string]interface{}{
		"id":       mailboxId,
		"email":    "update@test.com",
		"password": "new_password",
		"status":   1,
	}

	w4 := server.PUT(fmt.Sprintf("/api/user/mailboxes/%d", mailboxId), updateData, authHeader)
	if !server.AssertStatus(w4, http.StatusOK) {
		t.Fatalf("更新邮箱失败: %s", w4.Body.String())
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err := server.ParseResponse(w4, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("更新邮箱失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	t.Logf("✅ 更新邮箱测试通过")
}

// TestMailboxDelete 测试删除邮箱
func TestMailboxDelete(t *testing.T) {
	server := fixtures.NewTestServer()
	defer server.Close()

	// 插入测试数据
	server.SeedData()

	// 先注册用户并登录
	registerData := map[string]interface{}{
		"username": "deleteuser",
		"email":    "deleteuser@example.com",
		"password": "password123",
		"nickname": "删除测试用户",
	}

	server.POST("/api/public/user/register", registerData)

	loginData := map[string]interface{}{
		"username": "deleteuser",
		"password": "password123",
	}

	w2 := server.POST("/api/public/user/login", loginData)
	var loginResponse struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	server.ParseResponse(w2, &loginResponse)

	authHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", loginResponse.Data.Token),
	}

	// 创建邮箱
	mailboxData := map[string]interface{}{
		"domainId": 1, // 使用测试域名
		"email":    "delete@test.com",
		"password": "mailbox_password",
		"status":   1,
	}

	w3 := server.POST("/api/user/mailboxes", mailboxData, authHeader)
	var createResponse struct {
		Code int `json:"code"`
		Data struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}
	server.ParseResponse(w3, &createResponse)

	mailboxId := createResponse.Data.Id

	// 删除邮箱
	w4 := server.DELETE(fmt.Sprintf("/api/user/mailboxes/%d", mailboxId), authHeader)
	if !server.AssertStatus(w4, http.StatusOK) {
		t.Fatalf("删除邮箱失败: %s", w4.Body.String())
	}

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	err := server.ParseResponse(w4, &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if response.Code != 0 {
		t.Fatalf("删除邮箱失败: code=%d, msg=%s", response.Code, response.Msg)
	}

	t.Logf("✅ 删除邮箱测试通过")
}
