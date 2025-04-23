package application

import (
	"fmt"
	"net/http"
)

var (
	ErrDataNotFound        = &AppError{Code: http.StatusNotFound, Message: "数据不存在"}
	ErrInternalServerError = &AppError{Code: http.StatusInternalServerError, Message: "服务内部异常"}
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func NewAppError(code int, messageTemplate string, args ...any) *AppError {
	return &AppError{Code: code, Message: fmt.Sprintf(messageTemplate, args...)}
}

func (self *AppError) Error() string {
	return self.Message
}
