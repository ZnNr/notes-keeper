package handler

import (
	"github.com/ZnNr/notes-keeper.git/intenal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
