package app

import (
	"github.com/ZnNr/notes-keeper.git/intenal/config"
	router "github.com/ZnNr/notes-keeper.git/intenal/midleware/handler"
	"github.com/ZnNr/notes-keeper.git/intenal/midleware/repository"
	"github.com/ZnNr/notes-keeper.git/intenal/service"
	"github.com/ZnNr/notes-keeper.git/intenal/spellcheck"
	"github.com/ZnNr/notes-keeper.git/intenal/sqlite"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run(path string) {

	cfg := config.NewConfig(path)

	log := setupLogger(cfg.Level)
	log = log.With(slog.String("env", cfg.Level))
	log.Info("initializing server", slog.String("address", cfg.HTTP.Host+":"+cfg.HTTP.Port))
	log.Debug("logger debug mode enabled")

	// Создаём подключение к базе данных SQLite

	sqliteDB, err := sqlite.NewSQLiteConn(sqlite.SQLiteConfig{
		FilePath: cfg.SQLite.FilePath,
	})
	if err != nil {
		panic(err)
	}

	// Создаём репозитории с помощью sqliteDB
	repo := repository.NewRepositories(sqliteDB)
	//создаем таблицы
	repository.CreateTables(sqliteDB)

	// Создаём спеллчекер
	yaspeller := spellcheck.NewYaSpellChecker(log)

	// Создаём сервис с нужными зависимостями
	service := service.NewService(service.ServicesDependencies{
		Repos:    repo,
		Logger:   log,
		SignKey:  cfg.JWT.SignKey,  // Используем cfg.JWT.SignKey для JWT подписи
		TokenTTL: cfg.JWT.TokenTTL, // Используем cfg.JWT.TokenTTL для времени жизни токена
		Salt:     cfg.Hasher.Salt,  // Используем cfg.Hasher.Salt для хэширования
		Speller:  yaspeller,
	})

	// Логируем запуск сервера
	log.Info("starting server", slog.String("address", cfg.HTTP.Host+":"+cfg.HTTP.Port))

	// Создаём маршрутизатор
	r := chi.NewRouter()

	// Создаём обработчик с сервисом
	handler := router.NewHandler(service)

	// Создаём сервер с маршрутизатором и обработчиком
	server := router.NewServer(handler, r)

	// Инициализируем маршруты для сервера
	server.Router()

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
