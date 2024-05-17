package main

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	_ "github.com/mattn/go-sqlite3"
	"lab2/internal/handlers"
	"lab2/internal/repositories"
	"lab2/internal/services"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "./storage/lab2.db")

	if err != nil {
		log.Panicf("fail to connect db: %s", err.Error())
	}

	endpointRepo := repositories.NewEndpointRepo(db)
	endpointService := services.NewEndpointService(endpointRepo, endpointRepo)
	endpointHandler := handlers.NewEndpointHandler(endpointService)

	roleService := services.NewRoleService(repositories.NewRoleRepo(db), repositories.NewRoleRepo(db), endpointRepo)
	roleHandler := handlers.NewRoleHandler(roleService)

	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	userHandler := handlers.NewUserHandler(services.NewUserService(repositories.NewUserRepo(db), repositories.NewUserRepo(db)), tokenAuth)

	newHandler := handlers.NewNewHandler(services.NewNewService(repositories.NewNewCache(repositories.NewNewRepo(db))))

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Group(func(router chi.Router) {
		router.Get("/endpoints", endpointHandler.GetAll)

		router.Get("/roles", roleHandler.GetAll)
		router.Get("/roles/{name}", roleHandler.GetByID)
		router.Post("/roles", roleHandler.Create)
		router.Put("/roles/{name}", roleHandler.Update)
		router.Delete("/roles/{name}", roleHandler.Delete)

		router.Get("/users", userHandler.GetAll)
		router.Get("/users/{email}", userHandler.GetByID)
		router.Put("/users/{email}", userHandler.Update)

		router.Get("/news", newHandler.GetAll)
		router.Get("/news/{id}", newHandler.GetById)
	})

	router.Post("/signup", userHandler.SignUp)
	router.Post("/login", userHandler.Login)

	chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if _, err := endpointService.Create(ctx, method, route); err != nil {
			log.Printf("%v", err)
		}

		return nil
	})

	router.Handle("/storage/images/*", http.StripPrefix("/storage/images/", http.FileServer(http.Dir("storage/images/"))))

	http.ListenAndServe(":3333", router)
}
