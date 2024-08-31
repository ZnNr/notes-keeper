package handler

import (
	"encoding/json"
	"github.com/ZnNr/notes-keeper.git/intenal/errors"
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

type UserIDRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserIDResponse struct {
	Id int `json:"id"`
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

// GetUserIDHandler handles fetching the user ID
func (h *Handler) GetUserIDHandler(w http.ResponseWriter, r *http.Request) {
	var req UserIDRequest

	// Декодируем JSON-тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Вызываем метод GetUserID из AuthService
	userId, err := h.services.Auth.GetUserID(r.Context(), req.Username, req.Password)
	if err != nil {
		switch err {
		case errors.ErrUsernameRequired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.ErrPasswordRequired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.ErrCannotGetUser:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Формируем успешный ответ
	response := UserIDResponse{Id: userId}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
