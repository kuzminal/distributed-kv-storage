package router

import (
	"distapp/internal/app"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(i *app.Instance) http.Handler {
	r := chi.NewRouter()
	r.Mount("/debug", middleware.Profiler())
	r.Group(func(r chi.Router) {
		r.Post("/{id}", i.Set)
		r.Post("/notify", i.Notify)
		r.Get("/{id}", i.Get)
	})
	return r
}
