package services

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"lab2/internal/entities"
)

type UserReader interface {
	GetAll(ctx context.Context) ([]entities.User, error)
	GetByEmail(ctx context.Context, email string) (entities.User, error)
}

type UserWriter interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Update(ctx context.Context, user entities.User) (entities.User, error)
}

type UserService struct {
	r UserReader
	w UserWriter
}

func NewUserService(userReader UserReader, userWriter UserWriter) *UserService {
	return &UserService{r: userReader, w: userWriter}
}

func (s UserService) GetAll(ctx context.Context) ([]entities.User, error) {
	return s.r.GetAll(ctx)
}

func (s UserService) GetByEmail(ctx context.Context, email string) (entities.User, error) {
	return s.r.GetByEmail(ctx, email)
}

func (s UserService) GetByEmailAndPassword(ctx context.Context, email string, password string) (entities.User, error) {
	user, err := s.r.GetByEmail(ctx, email)

	if err != nil {
		return entities.User{}, err
	}

	if !s.checkPasswordHash(password, user.Password) {
		return entities.User{}, errors.New("bad pass")
	}

	return user, nil
}

func (s UserService) Create(ctx context.Context, user entities.User) (entities.User, error) {
	passwordBytes, err := s.hashPassword(user.Password)

	if err != nil {
		return entities.User{}, err
	}

	user.Role = "user"
	user.Password = passwordBytes

	return s.w.Create(ctx, user)
}

func (s UserService) Update(ctx context.Context, email string, user entities.User) (entities.User, error) {
	passwordBytes, err := s.hashPassword(user.Password)

	if err != nil {
		return entities.User{}, err
	}

	user.Password = passwordBytes

	return s.w.Update(ctx, user)
}

func (s UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
