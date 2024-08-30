package userdata

import (
	"context"
	"database/sql"
	"github.com/ZnNr/notes-keeper.git/intenal/errors"
	"github.com/ZnNr/notes-keeper.git/intenal/users/usermodel"
	_ "modernc.org/sqlite"
)

const (
	driverName = "sqlite"

	tableSchema = `
    CREATE TABLE  users (
    id serial PRIMARY KEY,
    username varchar(200),
    password varchar(200)
); 
`

	insertQuery = `
INSERT INTO users (username, password) VALUES (?, ?)
`
	getUserQuery = "SELECT id, username FROM users WHERE username = ? AND password = ?"

	checkUserExistsQuery = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CloseDb закрывает соединение с базой данных
func (data *UserRepository) CloseDb() {
	data.db.Close()
}

func (data *UserRepository) CreateUser(ctx context.Context, username, password string) (int64, error) {
	exists, err := data.checkUserExists(ctx, username)
	if err != nil {
		return 0, errors.ErrCannotCheckUserExist
	} else if exists {
		return 0, errors.ErrUserAlreadyExists
	}

	stmt, err := data.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, username, password)
	if err != nil {
		return 0, errors.ErrCannotCreateUser
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (r *UserRepository) checkUserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, checkUserExistsQuery, username).Scan(&exists)
	if err != nil {
		return false, errors.ErrCannotCheckUserExist
	}
	return exists, nil
}

func (r *UserRepository) GetUser(ctx context.Context, username, password string) (usermodel.User, error) {
	query := "SELECT id, username FROM users WHERE username = ? AND password = ?"
	row := r.db.QueryRowContext(ctx, query, username, password)
	var user usermodel.User
	err := row.Scan(&user.Id, &user.Username)
	return user, err
}
