package main

import (
	"fmt"
	"net/http"

	"github.com/Shercosta/digi-wallet/database"
	"github.com/Shercosta/digi-wallet/handlers"
	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	db := database.Connect()

	// r.Get("/", routes.GetBalance(db))
	// r.Get("/init-balance", routes.InitializeBalance(db))
	// r.Post("/take-balance", routes.PostTakeBalance(db))

	r.Post("/login", handlers.Login(db))
	r.Post("/register", handlers.Register(db))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/secure-route", func(w http.ResponseWriter, r *http.Request) {
			userID := middleware.GetUserID(r.Context())
			construct := map[string]any{
				"user_id": userID,
			}
			response.JSONSuccess(w, construct, nil, nil)
		})
	})

	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
