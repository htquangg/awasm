package cli

type AwasmResp[T any] struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
}
