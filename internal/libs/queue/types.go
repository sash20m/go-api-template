package queue

import "context"

const EventsExchangeName = "go-api-template"

const UsersCreatedQueueName = "users.created"
const UsersCreatedRoutingKey RoutingKey = "users.created"

// RabbitMQ Publisher types

type Message struct {
	Exchange    string
	RoutingKey  string
	Body        []byte
	ContentType string
}

type Publisher interface {
	Publish(ctx context.Context, msg Message) error
}

// RabbitMQ Consumer types

type RoutingKey string

type Binding struct {
	RoutingKey RoutingKey
}

type Delivery struct {
	Body       []byte
	RoutingKey RoutingKey
}

type RabbitMQHandler func(ctx context.Context, d Delivery) error

type Consumer interface {
	Consume(ctx context.Context, queueName string, handler RabbitMQHandler) error
}
