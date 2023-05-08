package auth

import (
	"strings"
	"time"

	"github.com/WebXense/micro"
	"github.com/WebXense/micro/jwt"
	"github.com/gin-gonic/gin"
)

const (
	ERR_CODE_UNAUTHORIZED     = "b97cf20d-42b6-470e-9e08-b4bb852c3811"
	ERR_CODE_FORBIDDEN        = "7792176d-0196-4a57-a959-93062c2b9b41"
	ERR_INTERNAL_SERVER_ERROR = "b6a82bc6-5884-41e1-8b6f-1a013b7da835"
	ERR_MSG_UNAUTHORIZED      = "Unauthorized"
	ERR_MSG_FORBIDDEN         = "Forbidden"
	ERR_MSG_INTERNAL_SERVER   = "Internal Server Error"
)

func init() {
	micro.RegisterError(ERR_CODE_UNAUTHORIZED, ERR_MSG_UNAUTHORIZED)
	micro.RegisterError(ERR_CODE_FORBIDDEN, ERR_MSG_FORBIDDEN)
	micro.RegisterError(ERR_INTERNAL_SERVER_ERROR, ERR_MSG_INTERNAL_SERVER)
}

func GetClaims(ctx *gin.Context) *jwt.Claims {
	token := GetAuthToken(ctx)
	if token == "" {
		return nil
	}

	var claims *jwt.Claims
	var err error
	// try to parse token with user token public key
	claims, err = jwt.ParseWithPublicKey(token, micro.USER_PUBLIC_PEM)
	if err != nil || claims == nil {
		// try to parse token with system token public key
		claims, err = jwt.ParseWithPublicKey(token, micro.SYSTEM_PUBLIC_PEM)
		if err != nil || claims == nil {
			return nil
		}
	}

	return claims
}

func GetAuthToken(ctx *gin.Context) string {
	t := ctx.GetHeader("Authorization")
	t = strings.ReplaceAll(t, "Bearer ", "")
	t = strings.ReplaceAll(t, " ", "")
	return t
}

func AbortUnauthorized(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	traces = append(traces, micro.Trace{
		Success:    false,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		SystemUUID: micro.SYSTEM_UUID,
		SystemName: micro.SYSTEM_NAME,
		TraceID:    traceID,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
	})
	var duration uint
	for _, t := range traces {
		duration += t.Duration
	}
	ctx.AbortWithStatusJSON(401, micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
		TraceID:  traceID,
		Traces:   traces,
		Duration: duration,
	})
}

func AbortForbidden(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	traces = append(traces, micro.Trace{
		Success:    false,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		SystemUUID: micro.SYSTEM_UUID,
		SystemName: micro.SYSTEM_NAME,
		TraceID:    traceID,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
	})
	var duration uint
	for _, t := range traces {
		duration += t.Duration
	}
	ctx.AbortWithStatusJSON(403, micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
		TraceID:  traceID,
		Traces:   traces,
		Duration: duration,
	})
}
