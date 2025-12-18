package controller

import "github.com/gin-gonic/gin"

type MissingPersonController interface {
	Create(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}