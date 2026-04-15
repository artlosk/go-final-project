package api

import (
	"go-final-project/internal/handlers"
	"go-final-project/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", handlers.NextHandler)
	r.Post("/api/signin", handlers.AuthHandler)

	r.Group(func(pr chi.Router) {
		pr.Use(middleware.Auth)
		pr.Post("/api/task", handlers.TaskHandler)
		pr.Get("/api/task", handlers.TaskHandler)
		pr.Put("/api/task", handlers.TaskHandler)
		pr.Delete("/api/task", handlers.TaskHandler)
		pr.Get("/api/tasks", handlers.TasksHandler)
		pr.Post("/api/task/done", handlers.TaskDoneHandler)
	})
}
