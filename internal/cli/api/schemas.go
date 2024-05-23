package api

const (
	UserAgent = "awasm-cli"
)

type AwasmResp[T any] struct {
	Data   T      `json:"data"`
	Reason string `json:"reason"`
	Msg    string `json:"msg"`
	Code   int    `json:"code"`
}
