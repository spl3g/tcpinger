package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
)

func NewServer(validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api/v1", apiRouter(validate))
	return r
}

func apiRouter(validate *validator.Validate) http.Handler {
	r := chi.NewRouter()
	r.Post("/check", handleCheck(validate))
	return r
}
