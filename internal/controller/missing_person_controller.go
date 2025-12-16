package controller

import "github.com/gin-gonic/gin"

type MissingPersonController interface {
	Create(ctx *gin.Context)
}