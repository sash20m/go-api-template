package queue

import "context"

type NoopPublisher struct{}

func (n NoopPublisher) Publish(ctx context.Context, msg Message) error {
	return nil
}

type NoopConsumer struct{}

func (n NoopConsumer) Consume(ctx context.Context, queueName string, handler RabbitMQHandler) error {
	<-ctx.Done()
	return ctx.Err()
}
