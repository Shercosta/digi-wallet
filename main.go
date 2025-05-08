package main

import (
	"fmt"
	"net/http"

	"github.com/Shercosta/digi-wallet/database"
	"github.com/Shercosta/digi-wallet/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	db := database.Connect()

	r.Get("/", routes.GetBalance(db))
	r.Get("/init-balance", routes.InitializeBalance(db))

	fmt.Printf("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
