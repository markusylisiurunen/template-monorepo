package main

import (
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/endpoints"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

func setupEventEndpoints(log logger.Logger) error {
	receiver, err := events.GetReceiver()
	if err != nil {
		return err
	}

	helloEndpoint, err := endpoints.NewEventEndpointHello(log)
	if err != nil {
		log.Errorw("failed to create 'hello' endpoint")
		return err
	}

	endpoints := []endpoints.EventEndpoint{
		helloEndpoint,
	}

	for _, endpoint := range endpoints {
		if err := endpoint.Register(receiver); err != nil {
			return err
		}
	}

	return nil
}
