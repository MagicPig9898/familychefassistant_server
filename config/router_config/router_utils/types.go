package router_utils

type GinResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(data any) GinResult {
	return GinResult{Code: 0, Msg: "ok", Data: data}
}

func Fail(msg string) GinResult {
	return GinResult{Code: -1, Msg: msg, Data: nil}
}
