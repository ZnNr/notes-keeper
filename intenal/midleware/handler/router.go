package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Server struct {
	handler *Handler
	router  *chi.Mux
	address string
}

func NewServer(handler *Handler, router *chi.Mux) *Server {
	return &Server{
		handler: handler,
		router:  router,
	}
}

func (s *Server) Router() {
	s.router.Use(middleware.Logger)

	s.router.Route("/note-keeper", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Post("/register", s.handler.RegisterHandler)
			r.Post("/login", s.handler.LoginHandler)
			r.Post("/get_user_id", s.handler.GetUserIDHandler)
			r.Group(func(r chi.Router) {

				r.Use(JWTMiddleware(s.handler.services.Auth))
				r.Post("/create", s.handler.CreateNoteHandler)
				r.Get("/notes", s.handler.GetNotesHandler)
			})
		})
	})

	http.ListenAndServe(":8080", s.router)
}
