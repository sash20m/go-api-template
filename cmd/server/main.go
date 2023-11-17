package main

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sash20m/go-api-template/config"
	api "github.com/sash20m/go-api-template/internal"
	"github.com/sash20m/go-api-template/pkg/logger"
	"github.com/sirupsen/logrus"
)

// @title Go Rest Api
// @description Api Endpoints for Go Server
func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"err": err,
		}).Error("Can't load config from .env. Problem with .env, or the server is in production environment.")
		return
	}

	config := config.ApiEnvConfig{
		Env:     strings.ToUpper(os.Getenv("ENV")),
		Port:    os.Getenv("PORT"),
		Version: os.Getenv("VERSION"),
	}

	logger.Log.WithFields(logrus.Fields{
		"env":     config.Env,
		"version": config.Version,
		"port":    config.Port,
	}).Info("Loaded app config")

	var wg sync.WaitGroup
	wg.Add(1)

	// Starting our magnificent server
	go func() {
		defer wg.Done()

		server := api.AppServer{}
		defer func() {
			if r := recover(); r != nil {
				server.OnShutdown()
			}
		}()

		server.Run(config)
	}()
	wg.Wait()

}

// cSpell:ignore logrus godotenv
