package api

import (
	"go-final-project/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", handlers.NextHandler)
	r.Post("/api/task", handlers.TaskHandler)
	r.Get("/api/task", handlers.TaskHandler)
	r.Put("/api/task", handlers.TaskHandler)
	r.Delete("/api/task", handlers.TaskHandler)
	r.Post("/api/task/done", handlers.TaskDoneHandler)
	r.Get("/api/tasks", handlers.TasksHandler)
}
