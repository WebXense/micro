package auth

import (
	"github.com/WebXense/micro/jwt"
	"github.com/gin-gonic/gin"
)

func SystemAdminOnly(ctx *gin.Context) {
	claims := GetClaims(ctx)
	if claims == nil {
		AbortUnauthorized(ctx)
		return
	}
	if claims.TokenType != jwt.TOKEN_TYPE_ADMIN_TOKEN && claims.TokenType != jwt.TOKEN_TYPE_SYSTEM_TOKEN {
		AbortUnauthorized(ctx)
		return
	}
}

func LoginRequired(ctx *gin.Context) {
	claims := GetClaims(ctx)
	if claims == nil {
		AbortUnauthorized(ctx)
		return
	}
	if claims.TokenType != jwt.TOKEN_TYPE_ACCESS_TOKEN {
		AbortUnauthorized(ctx)
		return
	}
	ctx.Next()
}

func RefreshTokenOnly(ctx *gin.Context) {
	claims := GetClaims(ctx)
	if claims == nil {
		AbortUnauthorized(ctx)
		return
	}
	if claims.TokenType != jwt.TOKEN_TYPE_REFRESH_TOKEN {
		AbortUnauthorized(ctx)
		return
	}
}
