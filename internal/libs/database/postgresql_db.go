package database

import (
	"errors"
	"fmt"

	"go-api-template/config"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // import file driver for migrations
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PostgresDB struct {
	Database *sqlx.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.CONFIG.DbHost,
		config.CONFIG.DbPort,
		config.CONFIG.DbUser,
		config.CONFIG.DbPassword,
		config.CONFIG.DbName,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{Database: db}, nil
}

// RunMigrations runs the migrations for the postgres db.
func (s *PostgresDB) RunMigrations(migrationsPath string) error {
	logrus.Info("Running migrations")
	pgStorage, err := postgres.WithInstance(s.Database.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Postgres instance with database failed: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", pgStorage)
	if err != nil {
		return fmt.Errorf("Migration with new database instance failed: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logrus.Info("No new migrations to run")
			return nil
		} else {
			return fmt.Errorf("m.Up err: %w", err)
		}
	}

	logrus.Info("Migrations script ran successfully")

	return nil
}

func (s *PostgresDB) Close() error {
	return s.Database.Close()
}
