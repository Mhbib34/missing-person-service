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
	"github.com/Mhbib34/missing-person-service/internal/model"
	"github.com/Mhbib34/missing-person-service/internal/repository"
	"github.com/Mhbib34/missing-person-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	err = db.AutoMigrate(&model.MissingPersons{})
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
		api.GET("/missing-persons/:id", controller.FindByID)
		api.GET("/missing-persons", controller.GetAll)
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


func TestGetMissingPersonByIdSuccess(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== create data via GORM (UUID auto) =====
	missingPerson := model.MissingPersons{
		Name:        "Joko",
		Age:         63,
		Description: "celana pendek",
		LastSeen:    "Medan",
		Contact:     "08123456789",
		PhotoID:     "test-image.jpg",
	}

	err := testDB.Create(&missingPerson).Error
	assert.Nil(t, err)

	// ⚠️ pastikan UUID ter-generate
	assert.NotEqual(t, uuid.Nil, missingPerson.ID)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/missing-persons/"+missingPerson.ID.String(),
		nil,
	)

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].(map[string]any)

	assert.Equal(t, missingPerson.ID.String(), data["id"])
	assert.Equal(t, "Joko", data["name"])
	assert.Equal(t, float64(63), data["age"])
	assert.Equal(t, "celana pendek", data["description"])
	assert.Equal(t, "Medan", data["last_seen"])
	assert.Equal(t, "08123456789", data["contact"])
	assert.Equal(t, "pending", data["image_status"])
	assert.Equal(t, "test-image.jpg", data["photo_id"])
}
func TestGetMissingPersonByIdFailedIfNotFound(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== create data via GORM (UUID auto) =====
	missingPerson := model.MissingPersons{
		Name:        "Joko",
		Age:         63,
		Description: "celana pendek",
		LastSeen:    "Medan",
		Contact:     "08123456789",
		PhotoID:     "test-image.jpg",
	}

	err := testDB.Create(&missingPerson).Error
	assert.Nil(t, err)

	// ⚠️ pastikan UUID ter-generate
	assert.NotEqual(t, uuid.Nil, missingPerson.ID)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/missing-persons/ef62bded-d467-4968-b686-742e256bd0b5",
		nil,
	)

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "NOT FOUND", response["status"])
	assert.Equal(t, "Report not found", response["error"])
}


func TestListMissingPersonSuccess(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== create data via GORM (UUID auto) =====
	missingPerson := model.MissingPersons{
		Name:        "Joko",
		Age:         63,
		Description: "celana pendek",
		LastSeen:    "Medan",
		Contact:     "08123456789",
		PhotoID:     "test-image.jpg",
		ImageStatus: "ready",
	}

	err := testDB.Create(&missingPerson).Error
	assert.Nil(t, err)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/missing-persons",
		nil,
	)

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 1)

	// ambil item pertama
	item := data[0].(map[string]any)

	assert.Equal(t, missingPerson.ID.String(), item["id"])
	assert.Equal(t, "Joko", item["name"])
	assert.Equal(t, float64(63), item["age"])
	assert.Equal(t, "celana pendek", item["description"])
	assert.Equal(t, "Medan", item["last_seen"])
	assert.Equal(t, "08123456789", item["contact"])
	assert.Equal(t, "ready", item["image_status"])
	assert.Equal(t, "test-image.jpg", item["photo_id"])
}

func TestListMissingPersonSuccessWithPagination(t *testing.T) {
	truncateMissingPersons(testDB)

	// ===== create data via GORM (UUID auto) =====
	missingPerson := model.MissingPersons{
		Name:        "Joko",
		Age:         63,
		Description: "celana pendek",
		LastSeen:    "Medan",
		Contact:     "08123456789",
		PhotoID:     "test-image.jpg",
		ImageStatus: "ready",
	}

	err := testDB.Create(&missingPerson).Error
	assert.Nil(t, err)

	// ===== request GET =====
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/missing-persons?limit=1&page=2",
		nil,
	)

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, req)

	// ===== assert response =====
	resp := recorder.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, _ := io.ReadAll(resp.Body)

	var response map[string]any
	_ = json.Unmarshal(respBody, &response)

	assert.Equal(t, "OK", response["status"])

	data := response["data"].([]any)
	assert.Len(t, data, 0)
}