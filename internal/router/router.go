package router

import (
	"github.com/Mhbib34/missing-person-service/internal/controller"
	"github.com/Mhbib34/missing-person-service/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller controller.MissingPersonController) *gin.Engine {
	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(middleware.ErrorRecovery()) // ⬅️ penting

	api := r.Group("/api/v1")
	{
		api.POST("/missing-persons", controller.Create)
	}

	return r
}
