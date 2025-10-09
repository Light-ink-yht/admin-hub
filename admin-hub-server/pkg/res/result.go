package res

// Result 返回响应
type Result struct {
	Code int    `json:"code"` // 响应状态码 0 成功，1失败
	Msg  string `json:"msg"`  // 响应消息
	Data any    `json:"data"` // 响应数据
}

// Success 成功响应，不带数据
func Success(msg string) Result {
	return Result{
		Code: 0, // 通常 0 表示成功
		Msg:  msg,
		Data: nil,
	}
}

// SuccessWithData 成功响应，带数据
func SuccessWithData(msg string, data interface{}) Result {
	return Result{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}

// FailWithWarn 失败警告响应，不带数据
func FailWithWarn(msg string) Result {
	return Result{
		Code: 1, // 警告
		Msg:  msg,
		Data: nil,
	}
}

// FailWithError 失败错误响应，不带数据
func FailWithError(msg string) Result {
	return Result{
		Code: 2, // 错误
		Msg:  msg,
		Data: nil,
	}
}

// NewResult 自定义响应
func NewResult(code int, msg string, data interface{}) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
