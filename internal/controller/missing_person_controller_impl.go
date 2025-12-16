package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mhbib34/missing-person-service/internal/dto"
	"github.com/Mhbib34/missing-person-service/internal/exception"
	"github.com/Mhbib34/missing-person-service/internal/helper"
	"github.com/Mhbib34/missing-person-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type MissingPersonControllerImpl struct {
	usecase usecase.MissingPersonUsecase
}

func NewMissingPersonController(u usecase.MissingPersonUsecase) MissingPersonController {
	return &MissingPersonControllerImpl{usecase: u}
}

func (c *MissingPersonControllerImpl) Create(ctx *gin.Context) {
	var request dto.CreateMissingPersonRequest

	if err := ctx.ShouldBind(&request); err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	result, err := c.usecase.Create(ctx.Request.Context(), request)
	if err != nil {
		exception.ErrorHandler(ctx, err)
		return
	}

	uploadDir := "storage/tmp"
	_ = os.MkdirAll(uploadDir, 0755)

	filename := request.Photo.Filename // atau uuid + ext
	filePath := filepath.Join(uploadDir, filename)

	ctx.SaveUploadedFile(request.Photo, filePath)

	webResponse := dto.WebResponse{
		Status: "OK",
		Message: "Report created successfully. Image is being processed.",
		Data:   result,
	}

	helper.WriteToResponseBody(ctx, http.StatusCreated, webResponse)
}