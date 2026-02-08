package queueTransport

import (
	"context"
	"go-api-template/internal/libs/queue"

	"github.com/sirupsen/logrus"
)

func HandleUsersCreated(ctx context.Context, d queue.Delivery) error {
	logrus.WithFields(logrus.Fields{
		"routingKey": d.RoutingKey,
		"bodySize":   len(d.Body),
	}).Info("Consumed message")

	// TODO: Handle the message in any way you want
	return nil
}
