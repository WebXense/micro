package micro

import (
	"encoding/json"
	"time"

	"github.com/WebXense/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron"
)

type Engine struct {
	GinEngine  *gin.Engine
	CronWorker *cron.Cron
	SystemUUID string
	SystemName string
}

func NewEngine(systemUUID, systemName string) *Engine {
	return &Engine{
		GinEngine:  gin.Default(),
		CronWorker: cron.New(),
		SystemUUID: systemUUID,
		SystemName: systemName,
	}
}

func (e *Engine) Run(addr string) {
	e.CronWorker.Start()
	e.GinEngine.GET("/micro/info", MicroInfoHandler(e))
	e.GinEngine.Run(addr)
}

func (e *Engine) RunServerOnly(addr string) {
	e.GinEngine.GET("/micro/info", MicroInfoHandler(e))
	e.GinEngine.Run(addr)
}

func (e *Engine) RunCronOnly() {
	e.CronWorker.Start()
}

func (e *Engine) Use(middleware ...gin.HandlerFunc) {
	e.GinEngine.Use(middleware...)
}

func GET[T any](engine *Engine, route string, handler Handler[T], middleware ...gin.HandlerFunc) {
	engine.GinEngine.GET(route, joinMiddlewareAndService(newGinServiceHandler(engine, handler), middleware...)...)
}

func GETWithCache[T any](engine *Engine, route string, cacheDuration time.Duration, handler Handler[T], middleware ...gin.HandlerFunc) {
	engine.GinEngine.GET(route, joinMiddlewareAndService(
		Cache(cacheDuration, newGinServiceHandler(engine, handler)), middleware...)...)
}

func POST[T any](engine *Engine, route string, handler Handler[T], middleware ...gin.HandlerFunc) {
	engine.GinEngine.POST(route, joinMiddlewareAndService(newGinServiceHandler(engine, handler), middleware...)...)
}

func PUT[T any](engine *Engine, route string, handler Handler[T], middleware ...gin.HandlerFunc) {
	engine.GinEngine.PUT(route, joinMiddlewareAndService(newGinServiceHandler(engine, handler), middleware...)...)
}

func DELETE[T any](engine *Engine, route string, handler Handler[T], middleware ...gin.HandlerFunc) {
	engine.GinEngine.DELETE(route, joinMiddlewareAndService(newGinServiceHandler(engine, handler), middleware...)...)
}

func Cron(engine *Engine, spec string, job func()) {
	engine.CronWorker.AddFunc(spec, job)
}

func newGinServiceHandler[T any](engine *Engine, handler Handler[T]) gin.HandlerFunc {
	handlerSetup := handler()
	return func(c *gin.Context) {
		traces := GetTraces(c)
		if len(traces) == 0 {
			traces = make([]Trace, 0)
		}
		traceID := GetTraceID(c)
		ctx := &Context[T]{
			GinContext: c,
			Request:    GinRequest[T](c),
			TraceID:    traceID,
		}
		if handlerSetup.Pagination {
			ctx.Page = GinRequest[sql.Pagination](c)
		}
		if handlerSetup.Sort {
			ctx.Sort = GinRequest[sql.Sort](c)
		}
		resp, err := handlerSetup.Service(ctx)
		if err != nil {
			traces = append(traces, Trace{
				TraceID:    traceID,
				Success:    false,
				Time:       time.Now().Format("2006-01-02 15:04:05"),
				SystemUUID: engine.SystemUUID,
				SystemName: engine.SystemName,
				Error: &ResponseError{
					Code:    err.Code(),
					Message: err.Error(),
				},
			})
			ctx.Error(err, traceID, traces)
			return
		}
		traces = append(traces, Trace{
			TraceID:    traceID,
			Success:    true,
			Time:       time.Now().Format("2006-01-02 15:04:05"),
			SystemUUID: engine.SystemUUID,
			SystemName: engine.SystemName,
		})
		ctx.OK(resp, traceID, traces, ctx.Page)
	}
}

func joinMiddlewareAndService(service gin.HandlerFunc, middleware ...gin.HandlerFunc) []gin.HandlerFunc {
	var funcs = make([]gin.HandlerFunc, 0)
	if len(middleware) > 0 {
		funcs = append(funcs, middleware...)
	}
	funcs = append(funcs, service)
	return funcs
}

func GetTraceID(c *gin.Context) string {
	traceID := c.GetHeader(MICRO_HEADER_TRACE_ID)
	if traceID == "" {
		return uuid.NewString()
	}
	return traceID
}

func SetTraceID(c *gin.Context, traceID string) {
	c.Request.Header.Set(MICRO_HEADER_TRACE_ID, traceID)
}

func GetTraces(c *gin.Context) []Trace {
	var traces []Trace
	traceHeader := c.GetHeader(MICRO_HEADER_TRACES)
	if traceHeader != "" {
		_ = json.Unmarshal([]byte(traceHeader), &traces)
	}
	if len(traces) == 0 {
		traces = make([]Trace, 0)
	}
	return traces
}

func SetTraces(c *gin.Context, traces []Trace) {
	b, _ := json.Marshal(traces)
	c.Request.Header.Set(MICRO_HEADER_TRACES, string(b))
}

type MicroInfo struct {
	SystemID   string `json:"system_id"`
	SystemName string `json:"system_name"`
}

func MicroInfoHandler(engine *Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, MicroInfo{
			SystemID:   engine.SystemUUID,
			SystemName: engine.SystemName,
		})
	}
}
