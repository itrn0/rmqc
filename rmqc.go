package rmqc

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var (
	ErrContextDone = errors.New("context done")
)

type Consumer struct {
	Url       string
	QueueName string
	options   *Options
	ctx       context.Context
	timeout   time.Duration
	messages  chan *Message
}

func NewConsumer(
	ctx context.Context,
	url string,
	queueName string,
	options *Options,
) *Consumer {
	return &Consumer{
		ctx:       ctx,
		Url:       url,
		QueueName: queueName,
		options:   handleOptions(options),
		messages:  make(chan *Message, options.Capacity),
	}
}

func (p *Consumer) Start() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}
		if err := p.connect(); err != nil {
			if errors.Is(err, ErrContextDone) {
				return
			}
			p.options.Logger.Warn("RabbitMQ connect error", "err", err)
		}
		p.options.Logger.Debug("reconnecting to RabbitMQ after 3 second")
		time.Sleep(3 * time.Second)
	}
}

func (p *Consumer) connect() error {
	conn, err := amqp.Dial(p.Url)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	msgCh, err := ch.Consume(
		p.QueueName,
		p.options.Consumer,
		p.options.AutoAck,
		p.options.Exclusive,
		p.options.NoLocal,
		p.options.NoWait,
		p.options.Args,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	p.options.Logger.Info("consumer connected to RabbitMQ", "queue", p.QueueName)
	defer p.options.Logger.Info("consumer disconnected from RabbitMQ", "queue", p.QueueName)

	for {
		select {
		case <-p.ctx.Done():
			return ErrContextDone
		case msg := <-msgCh:
			if p.options.Debug {
				p.options.Logger.Debug("received a message", "exchange", msg.Exchange, "routing_key",
					msg.RoutingKey, "delivery_tag", msg.DeliveryTag, "message_id", msg.MessageId)
			}
			if msg.DeliveryTag == 0 {
				if p.options.Debug {
					p.options.Logger.Debug("message delivery_tag is zero", "message_id", msg.MessageId)
				}
				continue
			}
			p.messages <- &Message{
				ExchangeName: msg.Exchange,
				RoutingKey:   msg.RoutingKey,
				Data:         msg.Body,
				ContentType:  msg.ContentType,
				Timestamp:    msg.Timestamp,
				MessageId:    msg.MessageId,
			}
			if p.options.Debug {
				p.options.Logger.Debug("message added to channel", "tag", msg.DeliveryTag,
					"message_id", msg.MessageId, "exchange", msg.Exchange, "routing_key", msg.RoutingKey)
			}
		}
	}
}

func (p *Consumer) Consume() <-chan *Message {
	return p.messages
}
