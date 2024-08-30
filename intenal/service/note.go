package service

import (
	"context"
	"github.com/ZnNr/notes-keeper.git/intenal/errors"
	"github.com/ZnNr/notes-keeper.git/intenal/midleware/repository"
	"github.com/ZnNr/notes-keeper.git/intenal/notes/notemodel"
	"github.com/ZnNr/notes-keeper.git/intenal/spellcheck"
	"log/slog"
	"strconv"
)

type NoteService struct {
	repo    repository.Note
	speller spellcheck.SpellChecker
	logger  *slog.Logger
}

func NewNoteService(noteRepo repository.Note, speller spellcheck.SpellChecker, logger *slog.Logger) *NoteService {
	return &NoteService{
		repo:    noteRepo,
		speller: speller,
		logger:  logger,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, userId int, text string) error {
	const op = "service.Note.CreateNote"
	s.logger = s.logger.With("op", op)

	if text == "" {
		return errors.ErrTextRequired
	}
	mistakes, err := s.speller.CheckText(text)
	if err != nil {
		s.logger.Error("cannot check mistakes", slog.String("userId", strconv.Itoa(userId)))
		return errors.ErrCannotCheckMistakes
	}
	_, err = s.repo.CreateNote(ctx, userId, text, mistakes)
	if err != nil {
		return err
	}
	return nil
}

func (s *NoteService) GetNotes(ctx context.Context, userId int) ([]notemodel.Note, error) {
	return s.repo.GetNotes(ctx, userId)
}
