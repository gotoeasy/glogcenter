package gweb

import (
	"encoding/json"

	"github.com/gotoeasy/glang/cmn"
)

type HttpResult struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success,omitempty"`
	Result  any    `json:"result,omitempty"`
}

func Error(code int, msg string) *HttpResult {
	return &HttpResult{
		Code:    code,
		Message: msg,
		Success: false,
	}
}

func Error500(msg string) *HttpResult {
	return &HttpResult{
		Code:    500,
		Message: msg,
		Success: false,
	}
}

func Error403() *HttpResult {
	return &HttpResult{
		Code:    403,
		Message: "forbidden",
		Success: false,
	}
}

func Error404() *HttpResult {
	return &HttpResult{
		Code:    404,
		Message: "page not found",
		Success: false,
	}
}

func Ok() *HttpResult {
	return &HttpResult{
		Code:    200,
		Success: true,
	}
}

func Ok200(msg string) *HttpResult {
	return &HttpResult{
		Code:    200,
		Message: msg,
		Success: true,
	}
}

func Result(result any) *HttpResult {
	return &HttpResult{
		Code:    200,
		Success: true,
		Result:  result,
	}
}

func (r *HttpResult) ToJson() string {
	bt, _ := json.Marshal(r)
	return cmn.BytesToString(bt)

}

func (r *HttpResult) LoadBytes(bytes []byte) error {
	return json.Unmarshal(bytes, r)
}
