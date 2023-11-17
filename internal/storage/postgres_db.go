package storage

import (
	"errors"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // import file driver for migrate
	"github.com/jmoiron/sqlx"
	"github.com/sash20m/go-api-template/pkg/logger"
)

func NewPostgresDB() (*Storage, error) {

	connStr := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " dbname=" + os.Getenv("DB_NAME") + " password=" + os.Getenv("DB_PASSWORD") + " sslmode=disable"
	logger.Log.Info(connStr, " eee")

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{db: db}, nil
}

// MigratePostgres migrates the postgres db to a new version.
func (s *Storage) MigratePostgres(migrationsPath string) error {
	pgStorage, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", pgStorage)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Info("No migrations to run")
		} else {
			return err
		}
	}

	logger.Log.Info("Migrations script ran successfully")

	return nil
}
