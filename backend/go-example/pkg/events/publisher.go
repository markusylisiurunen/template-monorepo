package events

import (
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

var publisherInstance *opinionatedevents.Publisher

func createPublisherInstance() error {
	destinations := []opinionatedevents.Destination{
		opinionatedevents.NewHTTPDestination("http://localhost:3081/events"),
	}

	publisher, err := opinionatedevents.NewPublisher(
		opinionatedevents.WithAsyncBridge(3, 500, destinations...),
	)

	if err != nil {
		return err
	}

	publisher.OnDeliveryFailure(func(msg *opinionatedevents.Message) {
		// TODO: persist failed messages
		logger.Default.Errorf("message %s failed to be delivered", msg.UUID())
	})

	publisherInstance = publisher

	return nil
}

func GetPublisher() (*opinionatedevents.Publisher, error) {
	if publisherInstance == nil {
		if err := createPublisherInstance(); err != nil {
			return nil, err
		}
	}

	return publisherInstance, nil
}
