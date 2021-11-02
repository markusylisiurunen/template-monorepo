package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/config"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/endpoints"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/middlewares"
)

func setupHttpEndpoints(cfg *config.Config, log logger.Logger) (*mux.Router, error) {
	router := mux.NewRouter()

	router.Use(
		handlers.CORS(
			handlers.AllowedMethods([]string{"HEAD", "OPTIONS", "GET", "PUT", "POST"}),
			handlers.AllowedOrigins(cfg.CorsAllowedOrigins),
			handlers.AllowedHeaders([]string{"Content-Type"}),
		),
		middlewares.JSON(),
	)

	// setup other middleware here

	sayHelloEndpoint, sayHelloEndpointError := endpoints.NewHttpEndpointSayHello(log)
	if sayHelloEndpointError != nil {
		log.Errorw("failed to create 'say hello' endpoint")
		return nil, sayHelloEndpointError
	}

	eventsEndpoint, eventsEndpointError := endpoints.NewHttpEndpointEvents(log)
	if eventsEndpointError != nil {
		log.Errorw("failed to create 'events' endpoint")
		return nil, eventsEndpointError
	}

	endpoints := []endpoints.HttpEndpoint{
		sayHelloEndpoint,
		eventsEndpoint,
	}

	for _, endpoint := range endpoints {
		if err := endpoint.Register(router); err != nil {
			return nil, err
		}
	}

	return router, nil
}
