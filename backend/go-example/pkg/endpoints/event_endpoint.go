package endpoints

import (
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
)

type EventEndpoint interface {
	Register(*opinionatedevents.Receiver) error
}
