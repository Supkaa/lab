package handlers

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"io"
	"lab2/internal/entities"
	"net/http"
)

type RoleService interface {
	GetAll(ctx context.Context) ([]entities.Role, error)
	GetByName(ctx context.Context, name string) (entities.Role, error)
	Create(ctx context.Context, name string, endpointIds []uuid.UUID) (entities.Role, error)
	Update(ctx context.Context, name string, endpointIds []uuid.UUID) (entities.Role, error)
	Delete(ctx context.Context, name string) error
}

type RoleHandler struct {
	rs RoleService
}

func NewRoleHandler(roleService RoleService) *RoleHandler {
	return &RoleHandler{rs: roleService}
}

type RoleResp struct {
	entities.Role
}

func (rr RoleResp) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewRoleListResp(roles []entities.Role) []render.Renderer {
	var list []render.Renderer

	for _, role := range roles {
		list = append(list, RoleResp{Role: role})
	}

	return list
}

func (h RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roles, err := h.rs.GetAll(ctx)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.RenderList(w, r, NewRoleListResp(roles)); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	role, err := h.rs.GetByName(ctx, chi.URLParam(r, "name"))

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.Render(w, r, RoleResp{Role: role}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string      `json:"name"`
		EndpointIds []uuid.UUID `json:"endpointIds"`
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

	if len(request.EndpointIds) == 0 {
		render.Render(w, r, BadRequestRender("endpointIds not be empty"))

		return
	}

	newRole, err := h.rs.Create(ctx, request.Name, request.EndpointIds)

	if err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}

	if err := render.Render(w, r, RoleResp{Role: newRole}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var request struct {
		EndpointIds []uuid.UUID `json:"endpointIds"`
	}

	ctx := r.Context()

	role, err := h.rs.GetByName(ctx, chi.URLParam(r, "name"))

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		if errors.Is(err, io.EOF) {
			render.Render(w, r, BadRequestRender("request body is empty"))

			return
		}

		render.Render(w, r, ErrRender(err))

		return
	}

	if len(request.EndpointIds) == 0 {
		render.Render(w, r, BadRequestRender("endpointIds not be empty"))

		return
	}

	updatedRole, err := h.rs.Update(ctx, role.Name, request.EndpointIds)

	if err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}

	if err := render.Render(w, r, RoleResp{Role: updatedRole}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.rs.Delete(ctx, chi.URLParam(r, "name")); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}

	if err := render.Render(w, r, &ErrResponse{HTTPStatusCode: 200}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}
