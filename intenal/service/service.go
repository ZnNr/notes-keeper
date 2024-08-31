package service

import (
	"context"
	"github.com/ZnNr/notes-keeper.git/intenal/midleware/repository"
	"github.com/ZnNr/notes-keeper.git/intenal/notes/notemodel"
	"github.com/ZnNr/notes-keeper.git/intenal/spellcheck"
	"log/slog"
	"time"
)

type Auth interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) error
	ParseToken(token string) (int, error)
	GetUserID(ctx context.Context, username string, password string) (int, error)
}

type Note interface {
	GetNotes(ctx context.Context, userId int) ([]notemodel.Note, error)
	CreateNote(ctx context.Context, userId int, text string) error
}

type Service struct {
	Auth
	Note
}

type ServicesDependencies struct {
	Repos    *repository.Repositories
	Logger   *slog.Logger
	SignKey  string
	TokenTTL time.Duration
	Salt     string
	Speller  spellcheck.SpellChecker
}

func NewService(deps ServicesDependencies) *Service {
	return &Service{
		Auth: NewAuthService(AuthDependencies{
			userRepo: deps.Repos.User,
			logger:   deps.Logger,
			signKey:  deps.SignKey,
			tokenTTL: deps.TokenTTL,
			salt:     deps.Salt,
		}),
		Note: NewNoteService(deps.Repos.Note, deps.Speller, deps.Logger),
	}
}
