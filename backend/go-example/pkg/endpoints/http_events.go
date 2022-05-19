package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type HttpEndpointEvents struct {
	log    logger.Logger
	handle opinionatedevents.ReceiveFromHTTP
}

func (endpoint *HttpEndpointEvents) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	endpoint.handle(resp, req)
}

func (endpoint *HttpEndpointEvents) Register(router *mux.Router) error {
	route := router.Path("/events").Methods(http.MethodPost)
	subrouter := route.Subrouter()

	// add middleware here if needed

	subrouter.NewRoute().HandlerFunc(endpoint.ServeHTTP)

	return nil
}

func NewHttpEndpointEvents(log logger.Logger) (*HttpEndpointEvents, error) {
	receiver, err := events.GetReceiver()
	if err != nil {
		return nil, err
	}

	return &HttpEndpointEvents{
		log:    log,
		handle: opinionatedevents.MakeReceiveFromHTTP(receiver),
	}, nil
}
