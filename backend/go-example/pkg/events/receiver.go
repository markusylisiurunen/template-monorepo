package events

import (
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
)

var receiverInstance *opinionatedevents.Receiver

func createReceiverInstance() error {
	receiver, err := opinionatedevents.NewReceiver()
	if err != nil {
		return err
	}

	receiverInstance = receiver

	return nil
}

func GetReceiver() (*opinionatedevents.Receiver, error) {
	if receiverInstance == nil {
		if err := createReceiverInstance(); err != nil {
			return nil, err
		}
	}

	return receiverInstance, nil
}
