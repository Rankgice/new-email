package result

// 通用错误码
var (
	ErrorReqParam     = ErrorResult(10000, "请求参数错误")
	ErrorBindingParam = ErrorResult(10001, "绑定参数错误")
	ErrorAdd          = ErrorResult(10100, "添加失败")
	ErrorUpdate       = ErrorResult(10101, "更新失败")
	ErrorDelete       = ErrorResult(10102, "删除失败")
	ErrorSelect       = ErrorResult(10103, "查询失败")
	ErrorCopy         = ErrorResult(10106, "复制失败")
	ErrorNotFound     = ErrorResult(10110, "未查询到数据")
	ErrorTxCommit     = ErrorResult(10113, "事务提交失败")
	ErrorDataVerify   = ErrorResult(10114, "数据校验失败")
)

// 认证相关错误码
var (
	ErrorUnauthorized   = ErrorResult(20001, "未授权访问")
	ErrorTokenInvalid   = ErrorResult(20002, "Token无效")
	ErrorTokenExpired   = ErrorResult(20003, "Token已过期")
	ErrorPermissionDeny = ErrorResult(20004, "权限不足")
	ErrorLoginFailed    = ErrorResult(20005, "登录失败")
	ErrorPasswordWrong  = ErrorResult(20006, "密码错误")
	ErrorUserNotFound   = ErrorResult(20007, "用户不存在")
	ErrorUserDisabled   = ErrorResult(20008, "用户已被禁用")
	ErrorAccountLocked  = ErrorResult(20009, "账户已被锁定")
)

// 邮件相关错误码
var (
	ErrorEmailSendFailed    = ErrorResult(30001, "邮件发送失败")
	ErrorEmailReceiveFailed = ErrorResult(30002, "邮件接收失败")
	ErrorEmailParseFailed   = ErrorResult(30003, "邮件解析失败")
	ErrorAttachmentTooLarge = ErrorResult(30004, "附件过大")
	ErrorAttachmentType     = ErrorResult(30005, "附件类型不支持")
	ErrorMailboxNotFound    = ErrorResult(30006, "邮箱不存在")
	ErrorMailboxDisabled    = ErrorResult(30007, "邮箱已被禁用")
	ErrorSMTPConfig         = ErrorResult(30008, "SMTP配置错误")
	ErrorIMAPConfig         = ErrorResult(30009, "IMAP配置错误")
)

// 规则相关错误码
var (
	ErrorRuleNotFound     = ErrorResult(40001, "规则不存在")
	ErrorRuleDisabled     = ErrorResult(40002, "规则已被禁用")
	ErrorRulePatternError = ErrorResult(40003, "规则模式错误")
	ErrorRuleConflict     = ErrorResult(40004, "规则冲突")
)

// 域名相关错误码
var (
	ErrorDomainNotFound    = ErrorResult(50001, "域名不存在")
	ErrorDomainNotVerified = ErrorResult(50002, "域名未验证")
	ErrorDNSVerifyFailed   = ErrorResult(50003, "DNS验证失败")
	ErrorDKIMGenFailed     = ErrorResult(50004, "DKIM生成失败")
)

// API相关错误码
var (
	ErrorAPIKeyInvalid  = ErrorResult(60001, "API密钥无效")
	ErrorAPIKeyExpired  = ErrorResult(60002, "API密钥已过期")
	ErrorAPIKeyDisabled = ErrorResult(60003, "API密钥已被禁用")
	ErrorAPIRateLimit   = ErrorResult(60004, "API调用频率超限")
	ErrorAPIPermission  = ErrorResult(60005, "API权限不足")
)

// 系统相关错误码
var (
	ErrorSystemMaintenance  = ErrorResult(70001, "系统维护中")
	ErrorSystemOverload     = ErrorResult(70002, "系统负载过高")
	ErrorConfigError        = ErrorResult(70003, "系统配置错误")
	ErrorServiceUnavailable = ErrorResult(70004, "服务不可用")
)

// ErrorResult 创建错误结果
func ErrorResult(code int, msg string) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}

// ErrorSimpleResult 创建简单错误结果
func ErrorSimpleResult(msg string) *Result {
	return &Result{
		Code: 20000,
		Msg:  msg,
	}
}
