package queueTransport

import (
	"context"
	"go-api-template/config"
	"go-api-template/internal/libs/queue"
	"go-api-template/internal/service"
)

type QueueTransport struct {
	services *service.Services
	rabbit   *queue.RabbitMQ
}

func NewQueueTransport(services *service.Services, rabbit *queue.RabbitMQ) *QueueTransport {
	return &QueueTransport{
		services: services,
		rabbit:   rabbit,
	}
}

func (t *QueueTransport) StartConsumers(ctx context.Context) error {
	if !config.CONFIG.RabbitMQEnabled || t.rabbit == nil {
		return nil
	}

	messagesRouter := NewRouter()
	messagesRouter.RegisterHandler(queue.UsersCreatedRoutingKey, HandleUsersCreated)

	type consumerDef struct {
		exchange string
		queue    string
	}

	consumers := []consumerDef{
		{
			exchange: queue.EventsExchangeName,
			queue:    queue.UsersCreatedQueueName,
		},
	}

	// Consume Queues
	errCh := make(chan error, len(consumers))
	for _, c := range consumers {
		go func(c consumerDef) {
			errCh <- t.rabbit.Consume(ctx, c.queue, messagesRouter.Handle)
		}(c)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}
