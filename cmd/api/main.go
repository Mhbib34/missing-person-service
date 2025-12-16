package main

import (
	"context"
	"log"
	"time"

	"github.com/Mhbib34/missing-person-service/internal/controller"
	"github.com/Mhbib34/missing-person-service/internal/database"
	"github.com/Mhbib34/missing-person-service/internal/repository"
	"github.com/Mhbib34/missing-person-service/internal/router"
	"github.com/Mhbib34/missing-person-service/internal/usecase"
	"github.com/Mhbib34/missing-person-service/internal/worker"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// lanjut connect DB
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewMissingPersonRepository(db)
	validate := validator.New()
	usecase := usecase.NewMissingPersonUsecase(repo, validate)
	controller := controller.NewMissingPersonController(usecase)

	ctx := context.Background()
	worker := worker.NewResizeImageJobWorker(db, 5)
	go worker.Start(ctx, 5*time.Second)


	r := router.SetupRouter(controller)
	r.Run(":3000")
}
