package amqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/streadway/amqp"

	"github.com/Bipolar-Penguin/svc-notificator/pkg/domain"
	"github.com/Bipolar-Penguin/svc-notificator/pkg/service"
)

const (
	reconnectTimeout = 5
	reinitTimeout    = 2
	promoCodeEntity  = "trading_session"
)

var (
	eventTypes        = [1]string{"update"}
	errChanNotCreated = errors.New("channel was not created")
)

type Consumer interface {
	Consume()
}

type Config struct {
	ExchangeName string
	ExchangeType string
	ConnString   string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	Args         amqp.Table

	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

type rabbitBroker struct {
	notificator     service.Notificator
	config          *Config
	connection      *amqp.Connection
	channel         *amqp.Channel
	logger          log.Logger
	eventsChan      chan<- domain.Event
	isBrokerActive  bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	quit            chan bool
}

func NewRabbitBroker(logger log.Logger, config *Config, notificator service.Notificator, eventsChan chan<- domain.Event) rabbitBroker {
	var rb = rabbitBroker{
		config:      config,
		logger:      logger,
		eventsChan:  eventsChan,
		notificator: notificator,
	}

	rb.quit = make(chan bool)

	return rb
}

func (rb rabbitBroker) Consume() {
	rb.handleReconnect()
}

func (rb *rabbitBroker) handleReconnect() {
	for {
		rb.isBrokerActive = false

		conn, err := rb.amqpConnect()
		if err != nil {
			rb.logger.Log("event", "retrying to connect")
			<-time.After(reconnectTimeout * time.Second)

			continue
		}

		rb.connection = conn
		rb.notifyConnClose = make(chan *amqp.Error)
		rb.connection.NotifyClose(rb.notifyConnClose)

		rb.handleReinit()
	}
}

func (rb *rabbitBroker) amqpConnect() (*amqp.Connection, error) {
	rb.logger.Log("event", "trying to connect via AMQP")

	conn, err := amqp.Dial(rb.config.ConnString)
	if err != nil {
		rb.logger.Log("error", err)
		return nil, err
	}

	rb.logger.Log("event", "AMQP connected")

	return conn, nil
}

func (rb *rabbitBroker) handleReinit() {
	for {
		rb.isBrokerActive = false

		ch, err := rb.initExchange()
		if err != nil {
			rb.logger.Log("error", err)
			<-time.After(reinitTimeout * time.Second)

			continue
		}

		rb.channel = ch
		rb.notifyChanClose = make(chan *amqp.Error)
		rb.channel.NotifyClose(rb.notifyChanClose)

		rb.isBrokerActive = true

		go rb.listenQueues()

		select {
		case <-rb.notifyConnClose:
			rb.logger.Log("event", "connection closed")
			rb.quit <- true

			return
		case <-rb.notifyChanClose:
			rb.logger.Log("event", "channel closed")
			rb.quit <- true
		}
	}
}

func (rb *rabbitBroker) initExchange() (*amqp.Channel, error) {
	rb.logger.Log("event", "trying to create an amqp channel")

	ch, err := rb.connection.Channel()
	if err != nil {
		rb.logger.Log("error", err)
		return nil, err
	}

	if ch == nil {
		return nil, errChanNotCreated
	}

	rb.logger.Log("event", "channel created")

	rb.logger.Log("event", "trying to create an AMQP exchange")

	err = ch.ExchangeDeclare(
		rb.config.ExchangeName, // exchange name
		rb.config.ExchangeType, // type
		rb.config.Durable,      // durable
		rb.config.AutoDelete,   // auto delete
		rb.config.Internal,     // internal
		rb.config.NoWait,       // no-wait
		rb.config.Args,         // arguments
	)
	if err != nil {
		rb.logger.Log("error", err)

		return nil, err
	}

	rb.logger.Log("event", "exchange created")

	return ch, nil
}

func (rb *rabbitBroker) listenQueues() {
	for _, eventType := range eventTypes {
		err := rb.declareQueue(eventType)
		if err != nil {
			rb.logger.Log("error", err)
		}

		go func(eventType string) {
			queueName := fmt.Sprintf("%s.%s", promoCodeEntity, eventType)

			rb.logger.Log("event", "started consuming", "queue", queueName)

			msgs, err := rb.channel.Consume(
				queueName,
				"",
				rb.config.AutoAck,
				rb.config.Exclusive,
				rb.config.NoLocal,
				rb.config.NoWait,
				rb.config.Args,
			)
			if err != nil {
				rb.logger.Log("error", err)
			}

			select {
			case <-msgs:
				for m := range msgs {
					var message domain.Event

					go rb.notificator.NotifyUser(&message)

					if err := json.Unmarshal(m.Body, &message); err != nil {
						rb.logger.Log("error", err)
					}

					rb.logger.Log("event", fmt.Sprintf("%v", message))

					rb.eventsChan <- message

					if err := m.Ack(false); err != nil {
						rb.logger.Log("error", err)
					}
				}
			case <-rb.quit:
				return
			}
		}(eventType)
	}
}

func (rb *rabbitBroker) declareQueue(event string) error {
	_, err := rb.channel.QueueDeclare(
		fmt.Sprintf("%s.%s", promoCodeEntity, event),
		rb.config.Durable,
		rb.config.AutoDelete,
		rb.config.Exclusive,
		rb.config.NoWait,
		rb.config.Args,
	)
	if err != nil {
		return err
	}

	return nil
}
