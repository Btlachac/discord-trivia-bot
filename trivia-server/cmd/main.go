package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"go-trivia-api/internal/db"
	"go-trivia-api/internal/server"
	"go-trivia-api/internal/service"
	"go-trivia-api/internal/worker"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-migrate/migrate/v4"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/sethvargo/go-envconfig"
)

type config struct {
	AudioFileDirectory string `env:"AUDIO_FILE_DIRECTORY,required"`

	DatabaseName   string `env:"DB_NAME,required"`
	DatabaseUser   string `env:"DB_USER,required"`
	DatabasePass   string `env:"DB_PASS,required"`
	DatabaseServer string `env:"DB_SERVER,required"`

	TriviaHostToken  string `env:"TRIVIA_HOST_TOKEN,required"`
	TriviaChannelID  string `env:"TRIVIA_CHANNEL_ID,required"`
	CommandChannelID string `env:"COMMAND_CHANNEL_ID,required"`
}

func main() {
	ctx := context.Background()

	var cfg config
	err := envconfig.Process(ctx, &cfg)
	check("envconfig.Process", err)

	router := mux.NewRouter().StrictSlash(true)
	db := getDBConnection(cfg)
	runMigrations(db)

	ts := createTriviaService(cfg, db)

	s := server.NewServer(router, ts)

	hostBot := getDiscordBot(cfg.TriviaHostToken)
	bot := worker.NewBot(
		hostBot,
		ts,
		cfg.TriviaChannelID,
		cfg.CommandChannelID,
		0,
		0,
		0,
		0,
	)

	eg := errgroup.Group{}

	eg.Go(func() error {
		return bot.Run()
	})
	eg.Go(func() error {
		return s.Run()
	})

	if err := eg.Wait(); err != nil {
		slog.Error("worker or server ended with error", err)
	} else {
		slog.Info("workers shut down successfully")
	}
}

func check(desc string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s -> %w", desc, err))
	}
}

func getDBConnection(cfg config) *sql.DB {
	dbConnStr := fmt.Sprintf("%s://%s:postgres@%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseServer,
		cfg.DatabaseUser)

	db, err := sql.Open("postgres", dbConnStr)

	check("sql.Open", err)

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := pgMigrate.WithInstance(db, &pgMigrate.Config{})
	check("pgMigrate.WithInstance", err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	check("migrate.NewWithDatabaseInstance", err)

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		check("m.Up", err)

	}
}

func createTriviaService(cfg config, sqlDB *sql.DB) *service.TriviaService {
	repository := db.NewTriviaRepository(sqlDB)
	service := service.NewTriviaService(repository, cfg.AudioFileDirectory)
	return service
}

func getDiscordBot(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	check("discordgo.New", err)
	return bot
}
