package http

func DefError(code int, msg string) *HttpResponse {
	return &HttpResponse{
		ErrCode: code,
		Message: msg,
	}
}

func DefSuccess(data interface{}) *HttpResponse {
	return &HttpResponse{
		Data: data,
	}
}

type HttpResponse struct {
	ErrCode int         `json:"err_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
