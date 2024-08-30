package handler

import (
	"encoding/json"
	"net/http"
)

type NoteReq struct {
	UserId int    `json:"user_id"`
	Text   string `json:"text"`
}

func (h *Handler) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userIdFromCtx := r.Context().Value("userID").(int)
	var note NoteReq

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	note.UserId = userIdFromCtx

	if err := h.services.Note.CreateNote(r.Context(), note.UserId, note.Text); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *Handler) GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userID").(int)
	notes, err := h.services.Note.GetNotes(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}
