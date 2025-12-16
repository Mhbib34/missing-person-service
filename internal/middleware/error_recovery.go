package middleware

import (
	"github.com/Mhbib34/missing-person-service/internal/exception"
	"github.com/gin-gonic/gin"
)

func ErrorRecovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				exception.ErrorHandler(ctx, err)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
