package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shercosta/digi-wallet/database"
	"github.com/Shercosta/digi-wallet/handlers"
	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/Shercosta/digi-wallet/routes"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	// It's good practice to handle the error, e.g., if .env isn't found
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, or error loading: %v. Relying on system environment variables.", err)
		// You can choose to log.Fatal here if .env is strictly required
	}
}

func main() {
	port := os.Getenv("PORT")
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

		r.Post("/take-balance", routes.PostTakeBalance(db))
		r.Get("/balance", routes.GetBalance(db))
	})

	// fmt.Println("Server running on port 3000")
	// http.ListenAndServe(":3000", r)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, r)
}
