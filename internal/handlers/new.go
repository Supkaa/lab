package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"lab2/internal/entities"
	"net/http"
	"slices"
	"strings"
)

type NewService interface {
	GetAll(ctx context.Context, orderBy string, order string) ([]entities.New, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.New, error)
}

type NewHandler struct {
	ns NewService
}

func NewNewHandler(newService NewService) *NewHandler {
	return &NewHandler{ns: newService}
}

type NewResp struct {
	entities.New
}

func (nr NewResp) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewNewListResp(news []entities.New) []render.Renderer {
	var list []render.Renderer

	for _, n := range news {
		list = append(list, NewResp{New: n})
	}

	return list
}

func (h NewHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderBy := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("orderBy")))
	order := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("order")))

	if orderBy == "" {
		orderBy = "created_at"
	}

	if order == "" {
		order = "ASC"
	}

	validOrderBy := []string{"created_at", "likes", "dislikes", "views"}
	validOrder := []string{"ASC", "DESC"}

	if !slices.Contains(validOrderBy, orderBy) {
		render.Render(w, r, BadRequestRender(fmt.Sprintf("orderBy must be equal %s", strings.Join(validOrderBy, ","))))
		return
	}

	if !slices.Contains(validOrder, strings.ToUpper(order)) {
		render.Render(w, r, BadRequestRender(fmt.Sprintf("order must be equal %s", strings.Join(validOrder, ","))))
		return
	}

	news, err := h.ns.GetAll(ctx, orderBy, order)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.RenderList(w, r, NewNewListResp(news)); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}

func (h NewHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))

	if err != nil {
		render.Render(w, r, BadRequestRender("wrong id type"))

		return
	}

	n, err := h.ns.GetByID(ctx, id)

	if err != nil {
		render.Render(w, r, &ErrResponse{
			Err:            err,
			HTTPStatusCode: http.StatusNotFound,
			ErrorText:      err.Error(),
		})

		return
	}

	if err := render.Render(w, r, NewResp{New: n}); err != nil {
		render.Render(w, r, ErrRender(err))

		return
	}
}
