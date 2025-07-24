package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shercosta/digi-wallet/middleware"
	"github.com/Shercosta/digi-wallet/models"
	"github.com/Shercosta/digi-wallet/request"
	"github.com/Shercosta/digi-wallet/response"
	"gorm.io/gorm"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response.JSONSuccess(w, "hello from routes", nil, nil)
}

func GetBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		userID := middleware.GetUserID(r.Context())

		fmt.Println("pass")

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}

func InitializeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		err := db.First(&balance).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				balance = models.Balance{Amount: 100000}
				if err := db.Create(&balance).Error; err != nil {
					response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
					return
				}
				response.JSONSuccess(w, "Balance initialized", nil, nil)
				return
			}

			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount = 100000
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, "Balance reset", nil, nil)
	}
}

func PostTakeBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var balance models.Balance

		var req request.TakeRequest
		req.AssignFormValues(r)

		userID := middleware.GetUserID(r.Context())

		err := db.Where("user_id = ?", userID).First(&balance).Error
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount -= *req.Amount
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, balance, nil, nil)
	}
}

func ListUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sort := r.URL.Query().Get("sort_amount")

		// Join user & balance
		query := db.Model(&models.User{}).Select("users.id, users.username, users.level, balances.amount").
			Joins("left join balances on users.id = balances.user_id")

		if sort == "asc" {
			query = query.Order("balances.amount asc")
		} else if sort == "desc" {
			query = query.Order("balances.amount desc")
		}

		type UserWithBalance struct {
			ID       uint    `json:"id"`
			Username string  `json:"username"`
			Level    int     `json:"level"`
			Amount   float64 `json:"amount"`
		}

		var result []UserWithBalance
		if err := query.Scan(&result).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, result, nil, nil)
	}
}

func PostAddBalance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Amount float64 `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSONError(w, http.StatusBadRequest, "invalid json", nil)
			return
		}

		userID := middleware.GetUserID(r.Context())

		var balance models.Balance
		if err := db.Where("user_id = ?", userID).First(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		balance.Amount += req.Amount
		if err := db.Save(&balance).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		// Update level based on balance
		var user models.User
		if err := db.First(&user, userID).Error; err == nil {
			switch {
			case balance.Amount > 3_000_000:
				user.Level = 4
			case balance.Amount > 2_000_000:
				user.Level = 3
			case balance.Amount > 1_000_000:
				user.Level = 2
			default:
				user.Level = 1
			}
			db.Save(&user)
		}

		response.JSONSuccess(w, map[string]any{
			"new_balance": balance.Amount,
			"new_level":   user.Level,
		}, nil, nil)
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil ID user yang login
		currentUserID := middleware.GetUserID(r.Context())

		// Ambil ID target dari URL
		idStr := r.URL.Path[len("/delete-user/"):]
		targetID, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}

		// Ambil user login & user target
		var currentUser, targetUser models.User
		if err := db.First(&currentUser, currentUserID).Error; err != nil {
			response.JSONError(w, http.StatusUnauthorized, "user not found", nil)
			return
		}

		if err := db.First(&targetUser, targetID).Error; err != nil {
			response.JSONError(w, http.StatusNotFound, "target user not found", nil)
			return
		}

		// Cek level
		if currentUser.Level <= targetUser.Level {
			response.JSONError(w, http.StatusForbidden, "you cannot delete a user with same or higher level", nil)
			return
		}

		// Hapus
		if err := db.Delete(&targetUser).Error; err != nil {
			response.JSONError(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		response.JSONSuccess(w, "user deleted", nil, nil)
	}
}
