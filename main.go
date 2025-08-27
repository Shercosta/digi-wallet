package main

import (
	"fmt"
	"net/http"

	"github.com/Shercosta/digi-wallet/database"
	"github.com/Shercosta/digi-wallet/handlers"
	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/Shercosta/digi-wallet/routes"
	"github.com/go-chi/chi/v5"
	"github.com/Shercosta/digi-wallet/database/migrations"
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
)

func main() {
	r := chi.NewRouter()

	db := database.Connect()


	// Public routes
	r.Post("/login", handlers.Login(db))
	r.Post("/register", handlers.Register(db))
	r.Get("/list-user", routes.ListUsers(db))

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddLevelToUsers(),
	})
	if err := m.Migrate(); err != nil {
		panic(err)
	}
	fmt.Println("Database migrated")

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Get("/secure-route", func(w http.ResponseWriter, r *http.Request) {
			userID := middleware.GetUserID(r.Context())
			construct := map[string]any{
				"user_id": userID,
			}
			response.JSONSuccess(w, construct, nil, nil)
		})

		//Balance
		r.Post("/take-balance", routes.PostTakeBalance(db))
		r.Get("/balance", routes.GetBalance(db))
		r.Put("/add-balance", routes.AddBalance(db))

		//User Management
		r.Delete("/delete-user/{id}", routes.DeleteUser(db))
	})

	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
