package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities/responses"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
)

type httpErrorResponseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type httpErrorResponse struct {
	Ok    bool                   `json:"ok"`
	Error httpErrorResponseError `json:"error"`
}

type httpDataResponse struct {
	Ok   bool               `json:"ok"`
	Data responses.Response `json:"data"`
}

type HttpRequest struct {
	*http.Request
}

type HttpResponseWriter struct {
	http.ResponseWriter
}

func (resp HttpResponseWriter) WriteError(statusCode int, httpError entities.HttpError) {
	body, err := json.Marshal(httpErrorResponse{
		Ok: false,
		Error: httpErrorResponseError{
			Code:    httpError.Code,
			Message: httpError.Message,
		},
	})

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(statusCode)

	if _, err := resp.Write(body); err != nil {
		logger.Default.Errorw("failed to write HTTP error",
			"Error", err.Error(),
		)
	}
}

func (resp HttpResponseWriter) WriteData(statusCode int, data responses.Response) {
	body, err := json.Marshal(httpDataResponse{
		Ok:   true,
		Data: data,
	})

	if err != nil {
		resp.WriteError(http.StatusInternalServerError, entities.HttpErrorUnknown)
		return
	}

	resp.WriteHeader(statusCode)

	if _, err := resp.Write(body); err != nil {
		logger.Default.Errorw("failed to write HTTP data",
			"Error", err.Error(),
		)
	}
}

type HttpEndpoint interface {
	Register(*mux.Router) error
}

type ServableHttpEndpoint interface {
	HttpEndpoint
	ServeHTTP(HttpResponseWriter, *HttpRequest)
}

func NewHttpEndpointHandler(endpoint ServableHttpEndpoint) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		rs := HttpResponseWriter{ResponseWriter: resp}
		rq := &HttpRequest{Request: req}

		endpoint.ServeHTTP(rs, rq)
	})
}
