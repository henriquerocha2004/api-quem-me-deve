package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	cliMemory "github.com/henriquerocha2004/quem-me-deve-api/client/memory"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	"github.com/henriquerocha2004/quem-me-deve-api/debt/gorm"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/container"
	"github.com/henriquerocha2004/quem-me-deve-api/internal/http/routes"
)

func main() {
	dependencies := fillDependencies()
	r := routes.Start(dependencies)
	http.Handle("/", r)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server started on port 8080")
}

func fillDependencies() *container.Dependencies {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	gormDB, err := gorm.NewGorm(dsn)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil
	}

	// debt dependencies
	debtRepo := gorm.NewGormDebtRepository(gormDB)
	cliRepo := cliMemory.NewClientDebtReader()

	service := debt.NewDebtService(debtRepo, cliRepo)

	return &container.Dependencies{
		DebtService: service,
	}
}
