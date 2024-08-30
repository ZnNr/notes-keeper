package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := h.services.Auth.Register(r.Context(), user.Username, user.Password); err != nil {
		// Добавляем логирование ошибки
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user LoginReq
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.services.Auth.Login(r.Context(), user.Username, user.Password)
	if err != nil {
		// Добавляем логирование ошибки
		http.Error(w, "Failed to authenticate user: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
