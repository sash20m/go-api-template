package main

import (
	"context"
	"go-api-template/config"
	"go-api-template/internal"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	serverConfig, err := config.LoadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"env":     serverConfig.Env,
		"version": serverConfig.Version,
		"port":    serverConfig.Port,
	}).Info("Loaded app config")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := internal.NewServer()
	defer func() {
		if r := recover(); r != nil {
			server.OnShutdown()
		}
	}()

	if err := server.Run(ctx); err != nil && ctx.Err() == nil {
		logrus.WithError(err).Fatal("Server exited with error")
	}
}
