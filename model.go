package micro

import "github.com/WebXense/sql"

type Response struct {
	Success    bool            `json:"success"`
	Error      *ResponseError  `json:"error,omitempty"`
	Pagination *sql.Pagination `json:"pagination,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
	Duration   uint            `json:"duration,omitempty"`
	TraceID    string          `json:"trace_id,omitempty"`
	Traces     []Trace         `json:"traces,omitempty"`
}

type ResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var errMap = make(map[interface{}]string)

func RegisterError(uuid string, message string) {
	errMap[uuid] = message
}

type Error interface {
	Code() string
	Error() string
}

func NewError(code string) Error {
	return &errorImp{
		code:    code,
		message: errMap[code],
	}
}

type errorImp struct {
	code    string
	message string
}

func (e *errorImp) Code() string {
	return e.code
}

func (e *errorImp) Error() string {
	return e.message
}
