package application

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	Message string                 `json:"message,omitempty"`
	Type    string                 `json:"type,omitempty"`
	Code    int                    `json:"code,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

const (
	DEFAULT_KEY         = "data"
	AUTH_ERROR          = "authError"
	NOT_FOUND_ERROR     = "notFoundError"
	BAD_PARAM_ERROR     = "badParamError"
	ALREADY_EXIST_ERROR = "alreadyExistError"
	NEW_SERVER_ERROR    = "newServerError"
	METHOD_NOT_ALLOWED  = "methodNotAllowed"
)

func (e *CustomError) Error() string {
	return fmt.Sprintf("Message:'%s', Type:%s, Code:%d, Data:%s", e.Message, e.Type, e.Code, e.Data)
}

func NewServerError(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusInternalServerError, Type: NEW_SERVER_ERROR}
}

func NewAlreadyExistsError(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusConflict, Type: ALREADY_EXIST_ERROR}
}

func NewNotFoundError(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusNotFound, Type: NOT_FOUND_ERROR}
}

func NewBadParamError(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusBadRequest, Type: BAD_PARAM_ERROR}
}

func NewAuthError(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusUnauthorized, Type: AUTH_ERROR}
}

func NewMethodNotAllowed(message string) *CustomError {
	return &CustomError{Message: message, Code: http.StatusMethodNotAllowed, Type: METHOD_NOT_ALLOWED}
}
