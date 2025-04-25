package handlers

import (
	"app/internal/models"
	"app/internal/services"
	"encoding/json"
	"net/http"
)

type authHandlers struct {
	svc services.AuthServiceInterface
}

func NewAuthHandlers(svc services.AuthServiceInterface) *authHandlers {
	return &authHandlers{
		svc: svc,
	}
}

func (h *authHandlers) SignUp(w http.ResponseWriter, r *http.Request) {
	var body *models.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Please provide correct input"))
		return
	}
	res, err := h.svc.SignUp(r.Context(), body)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, res.Status, res)
}

func (h *authHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var body *models.UserLoginParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Please provide correct input"))
		return
	}
	res, err := h.svc.Login(r.Context(), body)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, res.Status, res)
}

func (h *authHandlers) ConfirmAccount(w http.ResponseWriter, r *http.Request) {
	var body *models.UserConfirmationParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		models.ResponseWithJSON(w, http.StatusBadRequest, models.NewErrorResponse(http.StatusBadRequest, "Please provide correct input"))
		return
	}
	res, err := h.svc.ConfirmAccount(r.Context(), body)
	if err != nil {
		models.ResponseWithJSON(w, err.Status, err)
		return
	}
	models.ResponseWithJSON(w, res.Status, res)
}
