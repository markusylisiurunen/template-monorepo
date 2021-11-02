package endpoints

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities/messages"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type HttpEndpointEvents struct {
	log logger.Logger
}

func (endpoint *HttpEndpointEvents) ServeHTTP(resp HttpResponseWriter, req *HttpRequest) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, entities.HttpErrorInvalidPayload)
		return
	}

	msg, err := opinionatedevents.ParseMessage(body)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, entities.HttpErrorInvalidPayload)
		return
	}

	payload := &messages.HelloPayload{}
	if err := msg.Payload(payload); err != nil {
		resp.WriteError(http.StatusBadRequest, entities.HttpErrorInvalidPayload)
		return
	}

	endpoint.log.Infow("received a message",
		"UUID", msg.UUID(),
		"Greeting", payload.Greeting,
		"Name", payload.Name,
	)

	resp.WriteHeader(200)
}

func (endpoint *HttpEndpointEvents) Register(router *mux.Router) error {
	handler := NewHttpEndpointHandler(endpoint)

	route := router.Path("/events").Methods(http.MethodPost)
	subrouter := route.Subrouter()

	// add middleware here if needed

	subrouter.NewRoute().Handler(handler)

	return nil
}

func NewHttpEndpointEvents(log logger.Logger) (*HttpEndpointEvents, error) {
	return &HttpEndpointEvents{log: log}, nil
}
