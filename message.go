package rmqc

import (
	"fmt"
	"time"
)

type Message struct {
	Timestamp     time.Time
	Data          []byte
	ExchangeName  string
	RoutingKey    string
	ContentType   string
	MessageId     string
	ackFn         func(multiple bool) error
	nackFn        func(multiple bool, requeue bool) error
	ackedOrNacked bool
}

func (m *Message) ackNack() error {
	if m.ackedOrNacked {
		return fmt.Errorf("message already acked/nacked")
	}
	m.ackedOrNacked = true
	return nil
}

func (m *Message) Ack(multiple bool) error {
	if err := m.ackNack(); err != nil {
		return err
	}
	return m.ackFn(multiple)
}

func (m *Message) Nack(multiple bool, requeue bool) error {
	if err := m.ackNack(); err != nil {
		return err
	}
	return m.nackFn(multiple, requeue)
}
