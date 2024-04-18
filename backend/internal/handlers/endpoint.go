package handlers

import (
	"context"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"lab2/internal/entities"
	"net/http"
)

type EndpointService interface {
	GetAll(ctx context.Context) ([]entities.Endpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Endpoint, error)
	Create(ctx context.Context, method string, route string) (entities.Endpoint, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type EndpointHandler struct {
	es EndpointService
}

func NewEndpointHandler(endpointService EndpointService) *EndpointHandler {
	return &EndpointHandler{es: endpointService}
}

type EndpointResp struct {
	entities.Endpoint
}

func NewEndpointListResp(endpoints []entities.Endpoint) []render.Renderer {
	var list []render.Renderer

	for _, endpoint := range endpoints {
		list = append(list, EndpointResp{Endpoint: endpoint})
	}

	return list
}

func (er EndpointResp) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h EndpointHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endpoints, err := h.es.GetAll(ctx)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})
	}

	if err := render.RenderList(w, r, NewEndpointListResp(endpoints)); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	ErrorText string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		ErrorText:      err.Error(),
	}
}

func BadRequestRender(msg string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      msg,
	}
}
