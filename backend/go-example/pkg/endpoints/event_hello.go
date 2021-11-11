package endpoints

import (
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities/messages"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type EventEndpointHello struct {
	log logger.Logger
}

func (endpoint *EventEndpointHello) Handle(msg *opinionatedevents.Message) error {
	payload := &messages.HelloPayload{}
	if err := msg.Payload(payload); err != nil {
		return err
	}

	endpoint.log.Infow("received a message",
		"UUID", msg.UUID(),
		"Greeting", payload.Greeting,
		"Name", payload.Name,
	)

	return nil
}

func (endpoint *EventEndpointHello) Register(receiver *opinionatedevents.Receiver) error {
	return receiver.On("hello", endpoint.Handle)
}

func NewEventEndpointHello(log logger.Logger) (*EventEndpointHello, error) {
	return &EventEndpointHello{log: log}, nil
}
