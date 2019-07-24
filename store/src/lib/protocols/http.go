package protocols

func DefError(code int, msg string) *HttpResponse {
	return &HttpResponse{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

func DefSuccess(data interface{}) *HttpResponse {
	return &HttpResponse{
		Data: data,
	}
}

type HttpResponse struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}
