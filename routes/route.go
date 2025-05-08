package routes

import (
	"net/http"

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

		err := db.First(&balance).Error
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

		err := db.First(&balance).Error
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
