package swagger

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func CreateSwaggerRouter(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	return
}
