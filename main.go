package main

import (
	"fmt"
	"net/http"

	"github.com/Shercosta/digi-wallet/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", routes.HomeHandler)

	fmt.Printf("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
