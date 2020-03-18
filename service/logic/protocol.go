package logic

const (
	CodeSuccess = 10000 + iota
	CodeUnknownError
)

type HttpRequest struct {
}

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetResponse(code int, msg string, data interface{}) HttpResponse {
	return HttpResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func GetQuickResponse(data interface{}) HttpResponse {
	return HttpResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

func GetQuickErrorResponse(code int) HttpResponse {
	return HttpResponse{
		Code:    code,
		Message: "unknown error ",
		Data:    map[string]interface{}{},
	}
}
