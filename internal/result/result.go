package result

// Result 统一响应结构
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// SuccessResult 成功响应
func SuccessResult(data any) *Result {
	return NewResult(0, "success", data)
}

// ListResult 列表响应
func ListResult(data any, page, pageSize, total int64) *Result {
	return NewResult(0, "success", map[string]any{
		"list":     data,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// SimpleResult 简单成功响应
func SimpleResult(msg string) *Result {
	return NewResult(0, msg, nil)
}

// DataResult 带数据的成功响应
func DataResult(msg string, data any) *Result {
	return NewResult(0, msg, data)
}

// NewResult 创建响应结果
func NewResult(code int, msg string, data any) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// AddError 添加错误信息
func (r Result) AddError(err error) *Result {
	if err != nil {
		r.Msg += "：" + err.Error()
	}
	return &r
}
