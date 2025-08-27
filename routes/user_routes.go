package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/models"
	"github.com/Shercosta/digi-wallet/response"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func GetUserBalance(db *gorm.DB, userID uint) (*models.Balance, error) {
	var balance models.Balance
	err := db.Where("user_id = ?", userID).First(&balance).Error
	return &balance, err
}

func UpdateUserLevel(db *gorm.DB, user *models.User, balanceAmount float64) {
	switch {
	case balanceAmount > 3_000_000:
		user.Level = 4
	case balanceAmount > 2_000_000:
		user.Level = 3
	case balanceAmount > 1_000_000:
		user.Level = 2
	default:
		user.Level = 1
	}
	db.Save(user)
}

func ListUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		sort := r.URL.Query().Get("sort_amount")

		query := db.Model(&models.User{}).
			Joins("LEFT JOIN balances ON balances.user_id = users.id")

		if sort == "asc" {
			query = query.Order("COALESCE(balances.amount, 0) ASC")
		} else if sort == "desc" {
			query = query.Order("COALESCE(balances.amount, 0) DESC")
		}
		if err := query.Preload("Balance").Find(&users).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, users, nil, nil)
	}
}

func AddBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Amount float64 `json:"amount"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.JSONError(w, http.StatusBadRequest, "Invalid request", nil)
			return
		}

		userID := middleware.GetUserID(r.Context())
		balance, err := GetUserBalance(db, userID)
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount += body.Amount
		db.Save(balance)

		var user models.User
		db.First(&user, userID)
		UpdateUserLevel(db, &user, balance.Amount)

		response.JSONSuccess(w, map[string]interface{}{
			"balance": balance.Amount,
			"level":   user.Level,
		}, nil, nil)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetIDStr := chi.URLParam(r, "id")
		targetID, _ := strconv.Atoi(targetIDStr)

		currentID := middleware.GetUserID(r.Context())

		var currentUser, targetUser models.User
		if err := db.First(&currentUser, currentID).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		if err := db.First(&targetUser, targetID).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		if currentUser.Level <= targetUser.Level {
			response.JSONError(w, http.StatusForbidden, "Cannot delete user with equal or higher level", nil)
			return
		}

		db.Delete(&targetUser)
		response.JSONSuccess(w, "User deleted", nil, nil)
	}
}
