package controller

type errorResponse struct {
	StatusCode  int         `json:"status_code"`
	Msg         string      `json:"msg"`
	ErrorDetail interface{} `json:"error_detail"`
}
