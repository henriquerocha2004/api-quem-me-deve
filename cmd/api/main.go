package main

import (
	"fmt"
	"net/http"
	"time"

	cliMemory "github.com/henriquerocha2004/quem-me-deve-api/client/memory"
	"github.com/henriquerocha2004/quem-me-deve-api/debt"
	debtMemory "github.com/henriquerocha2004/quem-me-deve-api/debt/memory"
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

	// debt dependencies
	debtRepo := debtMemory.NewMemoryRepository()
	cliRepo := cliMemory.NewClientDebtReader()

	service := debt.NewDebtService(debtRepo, cliRepo)

	return &container.Dependencies{
		DebtService: service,
	}
}
