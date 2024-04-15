package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"io"
	"lab2/internal/entities"
	"log"
	"net/http"
	"time"
)

type UserService interface {
	GetAll(ctx context.Context) ([]entities.User, error)
	GetByEmail(ctx context.Context, email string) (entities.User, error)
	GetByEmailAndPassword(ctx context.Context, email string, password string) (entities.User, error)
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Update(ctx context.Context, email string, user entities.User) (entities.User, error)
}

type UserHandler struct {
	us  UserService
	jwt *jwtauth.JWTAuth
}

func NewUserHandler(userService UserService, jwt *jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{us: userService, jwt: jwt}
}

type UserResp struct {
	entities.User
	Token string `json:"token,omitempty"`
}

func (ur UserResp) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewUserListResp(users []entities.User) []render.Renderer {
	var list []render.Renderer

	for _, user := range users {
		list = append(list, UserResp{User: user})
	}

	return list
}

func (h UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, _ := jwtauth.FromContext(r.Context())

	log.Printf("%+v", claims)
	users, err := h.us.GetAll(ctx)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.RenderList(w, r, NewUserListResp(users)); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.us.GetByEmail(ctx, chi.URLParam(r, "email"))

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.Render(w, r, UserResp{User: user}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	ctx := r.Context()

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		if errors.Is(err, io.EOF) {
			render.Render(w, r, BadRequestRender("request body is empty"))

			return
		}

		render.Render(w, r, ErrRender(err))

		return
	}

	newUser, err := h.us.Create(ctx, entities.User{
		Name:     request.Name,
		Password: request.Password,
		Email:    request.Email,
	})

	if err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}

	if err := render.Render(w, r, UserResp{User: newUser}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	ctx := r.Context()

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		if errors.Is(err, io.EOF) {
			render.Render(w, r, BadRequestRender("request body is empty"))

			return
		}

		render.Render(w, r, ErrRender(err))

		return
	}

	user, err := h.us.GetByEmailAndPassword(ctx, request.Email, request.Password)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	claims := map[string]interface{}{"email": user.Email}
	jwtauth.SetExpiryIn(claims, 30*time.Minute)
	_, tokenString, err := h.jwt.Encode(claims)

	if err != nil {
		render.Render(w, r, BadRequestRender(err.Error()))

		return
	}

	if err := render.Render(w, r, UserResp{User: user, Token: tokenString}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	ctx := r.Context()

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		if errors.Is(err, io.EOF) {
			render.Render(w, r, BadRequestRender("request body is empty"))

			return
		}

		render.Render(w, r, ErrRender(err))

		return
	}

	newUser, err := h.us.Update(ctx, chi.URLParam(r, "email"), entities.User{
		Name:     request.Name,
		Password: request.Password,
		Email:    request.Email,
		Role:     request.Role,
	})

	if err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}

	if err := render.Render(w, r, UserResp{User: newUser}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}
