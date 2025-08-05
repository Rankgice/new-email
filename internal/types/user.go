package types

// UserProfileReq 用户资料更新请求
type UserProfileReq struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
}

// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Email    string `json:"email" binding:"required,email"`           // 邮箱
	Password string `json:"password" binding:"required,min=6"`        // 密码
	Nickname string `json:"nickname"`                                 // 昵称
	Code     string `json:"code"`                                     // 验证码（如果需要）
}

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"` // 密码
}

// UserLoginResp 用户登录响应
type UserLoginResp struct {
	Token string   `json:"token"` // 访问令牌
	User  UserResp `json:"user"`  // 用户信息
}

// UserSearchReq 用户搜索请求
type UserSearchReq struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"` // 搜索关键词
	Type    string `json:"type" form:"type"`                          // 搜索类型：username, email, nickname
	PageReq
}
