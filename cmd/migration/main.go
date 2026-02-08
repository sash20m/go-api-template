package main

import (
	"go-api-template/config"
	"go-api-template/internal/libs/database"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Loading config for migrations")
	_, err := config.LoadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config for migrations")
		panic(err)
	}

	postgresDB, err := database.NewPostgresDB()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create postgres db")
		panic(err)
	}

	logrus.Info("Connected to postgres db")

	// Print current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get current working directory")
		panic(err)
	}
	logrus.Infof("Current working directory: %s", cwd)
	err = postgresDB.RunMigrations("file://internal/migrations")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to run migrations")
		panic(err)
	}
}
