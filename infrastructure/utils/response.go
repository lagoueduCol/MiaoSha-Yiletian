package utils

type Response struct {
	Code int         `json:"code"` // 业务错误码
	Data interface{} `json:"data"` // 数据
	Msg  string      `json:"msg"`  // 提示信息
}
