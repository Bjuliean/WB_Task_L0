package response

import "net/http"

type Status int

const (
	OK            Status = http.StatusOK
	NotFound      Status = http.StatusNotFound
	BadRequest    Status = http.StatusBadRequest
	InternalError Status = http.StatusInternalServerError
)

type Response struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Send(st Status, msg string, err string) Response {
	return Response{
		Status:  st,
		Message: msg,
		Error:   err,
	}
}
