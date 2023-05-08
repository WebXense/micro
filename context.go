package micro

import (
	"github.com/WebXense/sql"
	"github.com/gin-gonic/gin"
)

type Trace struct {
	Success    bool           `json:"success"`
	Time       string         `json:"time"`
	SystemUUID string         `json:"system_uuid"`
	SystemName string         `json:"system_name"`
	TraceID    string         `json:"trace_id"`
	Duration   uint           `json:"duration"`
	Error      *ResponseError `json:"error,omitempty"`
}

type Context[T any] struct {
	GinContext *gin.Context
	TraceID    string
	Request    *T
	Page       *sql.Pagination
	Sort       *sql.Sort
	Response   interface{}
}

func (ctx *Context[T]) ClientIP() string {
	return ctx.GinContext.ClientIP()
}

func (ctx *Context[T]) UserAgent() string {
	return ctx.GinContext.Request.UserAgent()
}

func (ctx *Context[T]) OK(data interface{}, traceID string, traces []Trace, page ...*sql.Pagination) {
	var p *sql.Pagination
	if len(page) > 0 {
		p = page[0]
	}
	resp := &Response{
		Success:    true,
		Data:       data,
		Pagination: p,
		TraceID:    traceID,
		Traces:     traces,
	}
	ctx.Response = resp // for testing
	ctx.GinContext.JSON(200, resp)
}

func (ctx *Context[T]) Error(err Error, traceID string, traces []Trace) {
	resp := &Response{
		Success: false,
		Error: &ResponseError{
			Code:    err.Code(),
			Message: err.Error(),
		},
		TraceID: traceID,
		Traces:  traces,
	}
	ctx.Response = resp // for testing
	ctx.GinContext.JSON(200, resp)
}
