package endpoints

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities/responses"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type HttpEndpointSayHello struct {
	log logger.Logger
}

func (endpoint *HttpEndpointSayHello) ServeHTTP(resp HttpResponseWriter, req *HttpRequest) {
	params := mux.Vars(req.Request)
	name := params["name"]

	if strings.ToLower(name) == "donald" {
		resp.WriteError(http.StatusForbidden, entities.HttpErrorUnknown)
		return
	}

	resp.WriteData(200, responses.Message{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC(),
		Text:      fmt.Sprintf("Hello, %s 👋", name),
	})
}

func (endpoint *HttpEndpointSayHello) Register(router *mux.Router) error {
	handler := NewHttpEndpointHandler(endpoint)

	route := router.Path("/hello/{name}").Methods(http.MethodGet)
	subrouter := route.Subrouter()

	// add middleware here if needed

	subrouter.NewRoute().Handler(handler)

	return nil
}

func NewHttpEndpointSayHello(log logger.Logger) (*HttpEndpointSayHello, error) {
	return &HttpEndpointSayHello{
		log: log,
	}, nil
}
