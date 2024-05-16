package services

import (
	"context"
	"github.com/google/uuid"
	"lab2/internal/entities"
)

type NewReader interface {
	GetAll(ctx context.Context, orderBy string, order string) ([]entities.New, error)
	GetById(ctx context.Context, id uuid.UUID) (entities.New, error)
}

type NewService struct {
	r NewReader
}

func NewNewService(newReader NewReader) *NewService {
	return &NewService{r: newReader}
}

func (s NewService) GetAll(ctx context.Context, orderBy string, order string) ([]entities.New, error) {
	return s.r.GetAll(ctx, orderBy, order)
}

func (s NewService) GetByID(ctx context.Context, id uuid.UUID) (entities.New, error) {
	return s.r.GetById(ctx, id)
}
