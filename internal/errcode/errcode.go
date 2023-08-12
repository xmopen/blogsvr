// Package errcode  错误码.
package errcode

// CommonResult result struct.
type CommonResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) *CommonResult {
	return &CommonResult{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
}

var (
	// ErrorGetIndexArticleList 获取文章错误.
	ErrorGetIndexArticleList = &CommonResult{Code: 10010, Msg: "获取文章列表出错"}
)
