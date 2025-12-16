package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Mhbib34/missing-person-service/internal/controller"
	"github.com/Mhbib34/missing-person-service/internal/entity"
	"github.com/Mhbib34/missing-person-service/internal/middleware"
	"github.com/Mhbib34/missing-person-service/internal/repository"
	"github.com/Mhbib34/missing-person-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testRouter http.Handler
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=habib123 dbname=missing_person_test port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.MissingPersons{})
	if err != nil {
		panic(err)
	}

	return db
}

func setupRouter(db *gorm.DB) http.Handler {
	validate := validator.New()

	repo := repository.NewMissingPersonRepository(db)
	usecase := usecase.NewMissingPersonUsecase(repo, validate)
	controller := controller.NewMissingPersonController(usecase)

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

func truncateMissingPersons(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE missing_persons")
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testDB = setupTestDB()
	testRouter = setupRouter(testDB)

	code := m.Run()

	os.Exit(code)
}


func TestCreateMissingPersonSuccess(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== multipart body =====
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// text fields
	_ = writer.WriteField("name", "Joko")
	_ = writer.WriteField("age", "63")
	_ = writer.WriteField("description", "celana pendek")
	_ = writer.WriteField("last_seen", "Medan")
	_ = writer.WriteField("contact", "08123456789")

	// fake image file
	fileWriter, _ := writer.CreateFormFile(
		"photo",
		"test-image.jpg",
	)
	fileWriter.Write([]byte("FAKE_IMAGE_CONTENT"))

	writer.Close()

	// ===== request =====
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/missing-persons",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].(map[string]any)

	assert.Equal(t, "Joko", data["name"])
	assert.Equal(t, float64(63), data["age"])
	assert.Equal(t, "pending", data["image_status"])
	assert.Equal(t, "test-image.jpg", data["photo_id"])

	// ===== assert DB =====
	var count int64
	testDB.Model(&entity.MissingPersons{}).Count(&count)
	assert.Equal(t, int64(1), count)
}
func TestCreateMissingPersonFailedBadRequest(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== multipart body =====
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// text fields
	_ = writer.WriteField("name", "")
	_ = writer.WriteField("age", "63")
	_ = writer.WriteField("description", "celana pendek")
	_ = writer.WriteField("last_seen", "Medan")
	_ = writer.WriteField("contact", "08123456789")

	// fake image file
	fileWriter, _ := writer.CreateFormFile(
		"photo",
		"test-image.jpg",
	)
	fileWriter.Write([]byte("FAKE_IMAGE_CONTENT"))

	writer.Close()

	// ===== request =====
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/missing-persons",
		body,
	)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	json.Unmarshal(respBody, &response)

	assert.Equal(t, "BAD REQUEST", response["status"])
}
