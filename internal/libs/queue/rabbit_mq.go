package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	url string

	connMu sync.RWMutex
	conn   *amqp091.Connection

	pubMu sync.Mutex
	pubCh *amqp091.Channel

	consumeCh *amqp091.Channel

	prefetch int

	log *logrus.Entry
}

type RabbitMQOptions struct {
	Prefetch int
}

func NewRabbitMQ(url string, opts RabbitMQOptions) *RabbitMQ {
	prefetch := opts.Prefetch
	if prefetch <= 0 {
		prefetch = 20
	}
	return &RabbitMQ{
		url:      url,
		prefetch: prefetch,
		log:      logrus.WithField("component", "rabbitmq"),
	}
}

func (r *RabbitMQ) Connect() error {
	if r.url == "" {
		return errors.New("rabbitmq url is empty")
	}

	r.connMu.RLock()
	if r.conn != nil && !r.conn.IsClosed() && r.pubCh != nil && r.consumeCh != nil {
		r.connMu.RUnlock()
		return nil
	}
	r.connMu.RUnlock()

	r.connMu.Lock()
	defer r.connMu.Unlock()

	if r.conn != nil && !r.conn.IsClosed() && r.pubCh != nil && r.consumeCh != nil {
		return nil
	}

	conn, err := amqp091.Dial(r.url)
	if err != nil {
		return err
	}

	pubCh, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	consumeCh, err := conn.Channel()
	if err != nil {
		_ = pubCh.Close()
		_ = conn.Close()
		return err
	}

	if err := consumeCh.Qos(r.prefetch, 0, false); err != nil {
		_ = consumeCh.Close()
		_ = pubCh.Close()
		_ = conn.Close()
		return err
	}

	r.conn = conn
	r.pubCh = pubCh
	r.consumeCh = consumeCh

	r.log.WithFields(logrus.Fields{
		"prefetch": r.prefetch,
	}).Info("RabbitMQ connected")

	return nil
}

func (r *RabbitMQ) Close() error {
	r.connMu.Lock()
	defer r.connMu.Unlock()

	var closeErr error
	if r.consumeCh != nil {
		if err := r.consumeCh.Close(); err != nil {
			closeErr = err
		}
		r.consumeCh = nil
	}
	if r.pubCh != nil {
		if err := r.pubCh.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
		r.pubCh = nil
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
		r.conn = nil
	}
	return closeErr
}

func (r *RabbitMQ) Declare(exchange string, queueName string, bindings []Binding) error {
	if exchange == "" {
		return errors.New("exchange is empty")
	}
	if queueName == "" {
		return errors.New("queue name is empty")
	}

	r.connMu.RLock()
	ch := r.consumeCh
	r.connMu.RUnlock()
	if ch == nil {
		return errors.New("rabbitmq is not connected")
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	for _, b := range bindings {
		if b.RoutingKey == "" {
			continue
		}
		if err := ch.QueueBind(q.Name, string(b.RoutingKey), exchange, false, nil); err != nil {
			return err
		}
	}

	return nil
}

func (r *RabbitMQ) DeclareExchange(exchange string) error {
	if exchange == "" {
		return errors.New("exchange is empty")
	}

	r.connMu.RLock()
	ch := r.consumeCh
	r.connMu.RUnlock()
	if ch == nil {
		return errors.New("rabbitmq is not connected")
	}

	return ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
}

type QueueSpec struct {
	Name     string
	Bindings []Binding
}

// EnsureTopology ensures the exchange exists, then ensures each queue and its bindings exist.
// This is safe to call on every startup; declarations are idempotent as long as properties match.
func (r *RabbitMQ) EnsureTopology(exchange string, queues []QueueSpec) error {
	if err := r.DeclareExchange(exchange); err != nil {
		return err
	}
	for _, q := range queues {
		if q.Name == "" {
			continue
		}
		if err := r.Declare(exchange, q.Name, q.Bindings); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, msg Message) error {
	if msg.Exchange == "" {
		return errors.New("publish exchange is empty")
	}
	if msg.RoutingKey == "" {
		return errors.New("publish routing key is empty")
	}
	if msg.ContentType == "" {
		msg.ContentType = "application/json"
	}

	r.pubMu.Lock()
	defer r.pubMu.Unlock()

	r.connMu.RLock()
	ch := r.pubCh
	r.connMu.RUnlock()
	if ch == nil {
		return errors.New("rabbitmq not connected")
	}

	return ch.PublishWithContext(ctx, msg.Exchange, msg.RoutingKey, false, false, amqp091.Publishing{
		ContentType: msg.ContentType,
		Body:        msg.Body,
		Timestamp:   time.Now(),
	})
}

func (r *RabbitMQ) Consume(ctx context.Context, queueName string, handler RabbitMQHandler) error {
	if queueName == "" {
		return errors.New("queue name is empty")
	}
	if handler == nil {
		return errors.New("handler is nil")
	}

	r.connMu.RLock()
	ch := r.consumeCh
	r.connMu.RUnlock()
	if ch == nil {
		return errors.New("rabbitmq not connected")
	}

	consumerTag := fmt.Sprintf("pharos-%d", time.Now().UnixNano())
	deliveries, err := ch.Consume(queueName, consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	r.log.WithFields(logrus.Fields{
		"queue":       queueName,
		"consumerTag": consumerTag,
	}).Info("RabbitMQ consuming")

	for {
		select {
		case <-ctx.Done():
			_ = ch.Cancel(consumerTag, false)
			return ctx.Err()

		case d, ok := <-deliveries:
			if !ok {
				return errors.New("rabbitmq deliveries channel closed")
			}

			err := handler(ctx, Delivery{
				Body:       d.Body,
				RoutingKey: RoutingKey(d.RoutingKey),
			})
			if err != nil {
				d.Nack(false, true)
				continue
			}
			d.Ack(false)
		}
	}
}
