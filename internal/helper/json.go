package helper

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadFromRequestBody(r *http.Request, result any) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		panic(err)
	}
}

func WriteToResponseBody(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, data)
}