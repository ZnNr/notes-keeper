package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type SQLiteConfig struct {
	FilePath string
}

func NewSQLiteConn(config SQLiteConfig) (*sql.DB, error) {
	if config.FilePath == "" {
		return nil, fmt.Errorf("sqlite: file path is empty")
	}
	// Проверяем, существует ли файл базы данных, и создаем его, если отсутствует
	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		file, err := os.Create(config.FilePath)
		if err != nil {
			return nil, fmt.Errorf("sqlite: failed to create database file: %w", err)
		}
		file.Close() // Закрываем файл после создания
	}
	db, err := sql.Open("sqlite", config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("sqlite - New - sql.Open: %w", err)
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), defaultConnTimeout)
	defer cancel()
	if err := checkConnection(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func checkConnection(ctx context.Context, db *sql.DB) error {
	var lastErr error
	for i := 0; i < defaultConnAttempts; i++ {
		err := db.PingContext(ctx)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(defaultConnTimeout)
	}
	return fmt.Errorf("sqlite: can't connect to database after %d attempts, last error: %w", defaultConnAttempts, lastErr)
}
