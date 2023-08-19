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
	/*******通用错误码********/

	// ErrorParam 请求参数错误.
	ErrorParam = &CommonResult{Code: 10010, Msg: "请求参数错误"}

	// ErrorGetIndexArticleList 获取文章错误.
	ErrorGetIndexArticleList = &CommonResult{Code: 10010, Msg: "获取文章列表出错"}
)
