package repository

import (
	"context"
	"database/sql"
	"github.com/ZnNr/notes-keeper.git/intenal/notes/notedata"
	"github.com/ZnNr/notes-keeper.git/intenal/notes/notemodel"
	"github.com/ZnNr/notes-keeper.git/intenal/users/userdata"
	"github.com/ZnNr/notes-keeper.git/intenal/users/usermodel"
	_ "modernc.org/sqlite"
)

type User interface {
	CreateUser(ctx context.Context, username, password string) (int64, error)
	GetUser(ctx context.Context, username, password string) (usermodel.User, error)
}

type Note interface {
	CreateNote(ctx context.Context, userId int, text string, mistakes []byte) (int64, error)
	GetNotes(ctx context.Context, userId int) ([]notemodel.Note, error)
}
type Repositories struct {
	User
	Note
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: userdata.NewUserRepository(db),
		Note: notedata.NewNoteRepository(db),
	}
}
