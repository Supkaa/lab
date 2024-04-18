package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lab2/internal/entities"
)

type RoleReader interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, name string) (entities.Role, error)
}

type RoleWriter interface {
	Create(ctx context.Context, role entities.Role) (entities.Role, error)
	Update(ctx context.Context, name string, endpoints []entities.Endpoint) (entities.Role, error)
	Delete(ctx context.Context, name string) error
}

type RoleService struct {
	r  RoleReader
	w  RoleWriter
	er EndpointReader
}

func NewRoleService(roleReader RoleReader, roleWriter RoleWriter, endpointReader EndpointReader) *RoleService {
	return &RoleService{r: roleReader, w: roleWriter, er: endpointReader}
}

func (s RoleService) GetAll(ctx context.Context) ([]entities.Role, error) {
	return s.r.GetAll(ctx)
}

func (s RoleService) GetByName(ctx context.Context, name string) (entities.Role, error) {
	role, err := s.r.GetByName(ctx, name)

	if err != nil {
		return entities.Role{}, err
	}

	if role.Name == "" {
		return entities.Role{}, fmt.Errorf("role not found")
	}

	return role, nil
}

func (s RoleService) Create(ctx context.Context, name string, endpointIds []uuid.UUID) (entities.Role, error) {
	var endpoints []entities.Endpoint

	for _, endpointId := range endpointIds {
		endpoint, err := s.er.GetByID(ctx, endpointId)

		if err != nil {
			return entities.Role{}, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return s.w.Create(ctx, entities.Role{
		Name:      name,
		Endpoints: endpoints,
	})
}

func (s RoleService) Update(ctx context.Context, name string, endpointIds []uuid.UUID) (entities.Role, error) {
	var endpoints []entities.Endpoint

	for _, endpointId := range endpointIds {
		endpoint, err := s.er.GetByID(ctx, endpointId)

		if err != nil {
			return entities.Role{}, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return s.w.Update(ctx, name, endpoints)
}

func (s RoleService) Delete(ctx context.Context, name string) error {
	return s.w.Delete(ctx, name)
}
