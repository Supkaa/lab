package repositories

import (
	"context"
	"database/sql"
	"lab2/internal/entities"
)

type RoleRepo struct {
	db *sql.DB
}

func NewRoleRepo(db *sql.DB) RoleRepo {
	return RoleRepo{db: db}
}

func (r RoleRepo) GetAll(ctx context.Context) ([]entities.Role, error) {
	var roles []entities.Role
	rows, err := r.db.QueryContext(ctx,
		`SELECT r.*, p.* FROM roles r 
				INNER JOIN roles_have_endpoints rhe ON r.name = rhe.role_name 
    			INNER JOIN endpoints p ON rhe.endpoint_id = p.id
    			`)

	if err != nil {
		return nil, err
	}

	roleMap := make(map[string][]entities.Endpoint)

	for rows.Next() {
		var name string
		var endpoint entities.Endpoint
		if err := rows.Scan(&name, &endpoint.ID, &endpoint.URL, &endpoint.Method); err != nil {
			return nil, err
		}

		roleMap[name] = append(roleMap[name], endpoint)
	}

	for role, endpoints := range roleMap {
		roles = append(roles, entities.Role{
			Name:      role,
			Endpoints: endpoints,
		})
	}

	return roles, nil
}

func (r RoleRepo) GetByName(ctx context.Context, roleName string) (entities.Role, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT r.*, p.* FROM roles r 
				INNER JOIN roles_have_endpoints rhe ON r.name = rhe.role_name 
    			INNER JOIN endpoints p ON rhe.endpoint_id = p.id
				WHERE r.name = $1
    			`, roleName)

	if err != nil {
		return entities.Role{}, err
	}

	var name string
	var endpoints []entities.Endpoint

	for rows.Next() {
		var endpoint entities.Endpoint

		if err := rows.Scan(&name, &endpoint.ID, &endpoint.URL, &endpoint.Method); err != nil {
			return entities.Role{}, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return entities.Role{
		Name:      name,
		Endpoints: endpoints,
	}, nil
}

func (r RoleRepo) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return entities.Role{}, err
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO roles(name) VALUES (?)", role.Name); err != nil {
		if err := tx.Rollback(); err != nil {
			return entities.Role{}, err
		}

		return entities.Role{}, err
	}

	for _, endpoint := range role.Endpoints {
		if _, err := tx.ExecContext(ctx, "INSERT INTO roles_have_endpoints(role_name, endpoint_id) VALUES (?, ?)", role.Name, endpoint.ID); err != nil {
			if err := tx.Rollback(); err != nil {
				return entities.Role{}, err
			}

			return entities.Role{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return entities.Role{}, err
	}

	return role, nil
}

func (r RoleRepo) Update(ctx context.Context, name string, endpoints []entities.Endpoint) (entities.Role, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return entities.Role{}, err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE roles SET name = $1 WHERE name = $1", name); err != nil {
		if err := tx.Rollback(); err != nil {
			return entities.Role{}, err
		}

		return entities.Role{}, err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM roles_have_endpoints WHERE role_name = ?", name); err != nil {
		if err := tx.Rollback(); err != nil {
			return entities.Role{}, err
		}

		return entities.Role{}, err
	}

	for _, endpoint := range endpoints {
		if _, err := tx.ExecContext(ctx, "INSERT INTO roles_have_endpoints(role_name, endpoint_id) VALUES (?, ?)", name, endpoint.ID); err != nil {
			if err := tx.Rollback(); err != nil {
				return entities.Role{}, err
			}

			return entities.Role{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return entities.Role{}, err
	}

	return entities.Role{
		Name:      name,
		Endpoints: endpoints,
	}, nil
}

func (r RoleRepo) Delete(ctx context.Context, name string) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM roles WHERE name = ?", name); err != nil {
		return err
	}

	return nil
}
