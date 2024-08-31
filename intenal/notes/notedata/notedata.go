package notedata

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/ZnNr/notes-keeper.git/intenal/errors"
	"github.com/ZnNr/notes-keeper.git/intenal/notes/notemodel"
	_ "modernc.org/sqlite"
)

const (
	driverName = "sqlite"

	insertQuery = `
INSERT INTO notes (user_id, text, mistakes) VALUES (?, ?, ?)
`
	getNoteQuery = "SELECT * FROM notes WHERE id = ?"

	getNotesQuery = "SELECT * FROM notes WHERE user_id = ?"

	deleteQuery = "DELETE FROM notes WHERE id=:id"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

// CreateNote создает новую заметку в базе данных и возвращает её ID.
func (data *NoteRepository) CreateNote(ctx context.Context, userId int, text string, mistakes []byte) (int64, error) {
	stmt, err := data.db.PrepareContext(ctx, insertQuery)
	if err != nil {
		return 0, errors.ErrCannotPrepareStatement
	}
	defer stmt.Close()
	if err != nil {
		return 0, errors.ErrCannotCreateNote
	}
	res, err := stmt.ExecContext(ctx, userId, text, mistakes)
	if err != nil {
		return 0, errors.ErrCannotCreateNote
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// CloseDb закрывает соединение с базой данных
func (data *NoteRepository) CloseDb() {
	data.db.Close()
}

func (data NoteRepository) GetNote(id int) (notemodel.Note, error) {

	row := data.db.QueryRow(getNoteQuery, id)

	var note notemodel.Note
	err := row.Scan(&note.Id, &note.UserId, &note.Text, &note.Mistakes)
	return note, err
}

// GetNotes извлекает заметки пользователя из базы данных.
func (r *NoteRepository) GetNotes(ctx context.Context, userId int) ([]notemodel.Note, error) {
	rows, err := r.db.QueryContext(ctx, getNotesQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []notemodel.Note
	for rows.Next() {
		var note notemodel.Note
		var mistakesBytes []byte // временная переменная для хранения JSON-данных

		// Сканируем данные
		err := rows.Scan(&note.Id, &note.UserId, &note.Text, &mistakesBytes)
		if err != nil {
			return nil, err
		}

		// Преобразуем JSON-данные в структуру Go
		if len(mistakesBytes) > 0 {
			err = json.Unmarshal(mistakesBytes, &note.Mistakes)
			if err != nil {
				return nil, err
			}
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (data *NoteRepository) Delete(id int) (bool, error) {
	// Получаем заметку по ID для проверки существования
	_, err := data.GetNote(id)
	if err != nil {
		return false, err
	}

	res, err := data.db.Exec(deleteQuery, sql.Named("id", id))
	if err != nil {
		return false, err
	}

	deleted, err := res.RowsAffected()
	return deleted == 1, err
}
