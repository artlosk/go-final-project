package handlers

import (
	"encoding/json"
	"go-final-project/internal/auth"
	"go-final-project/internal/helpers"
	"net/http"
	"os"
	"strings"
)

type signInRequest struct {
	Password string `json:"password"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	serverPass := os.Getenv("TODO_PASSWORD")
	if serverPass == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "пароль сервера не задан"})
		return
	}
	if strings.TrimSpace(req.Password) == "" || req.Password != serverPass {
		helpers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "неверный пароль"})
		return
	}
	token, err := auth.MakeToken(serverPass)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
