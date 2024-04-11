package repositories

import (
	"context"
	"database/sql"
	"lab2/internal/entities"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r UserRepo) GetAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	rows, err := r.db.QueryContext(ctx, "SELECT email, name, password, role FROM users")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user entities.User

		if err := rows.Scan(&user.Email, &user.Name, &user.Password, &user.Role); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r UserRepo) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	row := r.db.QueryRowContext(ctx, "SELECT email, name, password, role FROM users WHERE email = $1", email)

	if err := row.Scan(&user.Email, &user.Name, &user.Password, &user.Role); err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users(email, name, password, role) VALUES (?, ?, ?, ?)", user.Email, user.Name, user.Password, user.Role); err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r UserRepo) Update(ctx context.Context, user entities.User) (entities.User, error) {
	if _, err := r.db.ExecContext(ctx, "UPDATE users SET email = $1, name = $2, password = $3, role = $4 WHERE email = $1", user.Email, user.Name, user.Password, user.Role); err != nil {
		return entities.User{}, err
	}

	return user, nil
}
