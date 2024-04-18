package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lab2/internal/entities"
)

type EndpointReader interface {
	GetAll(ctx context.Context) ([]entities.Endpoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.Endpoint, error)
}

type EndpointWriter interface {
	Create(ctx context.Context, endpoint entities.Endpoint) (entities.Endpoint, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type EndpointService struct {
	r EndpointReader
	w EndpointWriter
}

func NewEndpointService(endpointReader EndpointReader, endpointWriter EndpointWriter) *EndpointService {
	return &EndpointService{r: endpointReader, w: endpointWriter}
}

func (s EndpointService) GetAll(ctx context.Context) ([]entities.Endpoint, error) {
	return s.r.GetAll(ctx)
}

func (s EndpointService) GetByID(ctx context.Context, id uuid.UUID) (entities.Endpoint, error) {
	return s.r.GetByID(ctx, id)
}

func (s EndpointService) Create(ctx context.Context, method string, route string) (entities.Endpoint, error) {
	baseId := uuid.MustParse("1d03b612-d864-49c3-90bd-0ea2f5114d07")

	uuid.NewSHA1(baseId, []byte(fmt.Sprintf("%s%s", method, route)))
	return s.w.Create(ctx, entities.Endpoint{
		ID:     uuid.NewSHA1(baseId, []byte(fmt.Sprintf("%s%s", method, route))),
		URL:    route,
		Method: method,
	})
}

func (s EndpointService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.w.Delete(ctx, id)
}
