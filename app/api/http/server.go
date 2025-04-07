package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (s *Task) WithTaskHandlers(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/", func(r chi.Router) {
		r.Get("/tasks", s.allTasks)
		r.With(authMiddleware).Get("/status/{id}", s.taskStatus)
		r.With(authMiddleware).Get("/result/{id}", s.taskResult)
		r.With(authMiddleware).Post("/task", s.newTask)
	})
}

func (u *User) WithAuthHandlers(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", u.registerHandler)
		r.Post("/login", u.loginHandler)
	})
}
