package logic

const (
	CodeSuccess = iota
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
		Code: code,
		Message: msg,
		Data: data,
	}
}

func GetQuickResponse(data interface{}) HttpResponse {
	return HttpResponse{
		Code: CodeSuccess,
		Message: "",
		Data: data,
	}
}

func (s *Service) Response() {

}
func (s *Service) Request() {

}
