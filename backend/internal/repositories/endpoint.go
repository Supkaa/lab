package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"lab2/internal/entities"
)

type EndpointRepo struct {
	db *sql.DB
}

func NewEndpointRepo(db *sql.DB) EndpointRepo {
	return EndpointRepo{db: db}
}

func (r EndpointRepo) GetAll(ctx context.Context) ([]entities.Endpoint, error) {
	var endpoints []entities.Endpoint
	rows, err := r.db.QueryContext(ctx, "SELECT id, url, method FROM endpoints")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var endpoint entities.Endpoint

		if err := rows.Scan(&endpoint.ID, &endpoint.URL, &endpoint.Method); err != nil {
			return nil, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func (r EndpointRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.Endpoint, error) {
	var endpoint entities.Endpoint

	rows, err := r.db.QueryContext(ctx, "SELECT id, url, method FROM endpoints WHERE id = $1", id)

	if err != nil {
		return entities.Endpoint{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&endpoint.ID, &endpoint.URL, &endpoint.Method); err != nil {
			return entities.Endpoint{}, err
		}
	}

	return endpoint, nil
}

func (r EndpointRepo) Create(ctx context.Context, endpoint entities.Endpoint) (entities.Endpoint, error) {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO endpoints(id, url, method) VALUES ($1, $2, $3)", endpoint.ID, endpoint.URL, endpoint.Method); err != nil {
		return entities.Endpoint{}, err
	}

	return endpoint, nil
}

func (r EndpointRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM endpoints WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
