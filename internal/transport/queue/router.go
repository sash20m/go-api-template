package queueTransport

import (
	"context"
	"go-api-template/internal/libs/queue"
)

type Router struct {
	handlers map[queue.RoutingKey]queue.RabbitMQHandler
}

func NewRouter() *Router {
	return &Router{
		handlers: map[queue.RoutingKey]queue.RabbitMQHandler{},
	}
}

func (r *Router) RegisterHandler(routingKey queue.RoutingKey, h queue.RabbitMQHandler) {
	if routingKey == "" || h == nil {
		return
	}
	r.handlers[routingKey] = h
}

func (r *Router) Handle(ctx context.Context, d queue.Delivery) error {
	h, ok := r.handlers[d.RoutingKey]
	if !ok {
		// Unknown routing key: ack/drop to avoid infinite requeue loops.
		return nil
	}
	return h(ctx, d)
}
