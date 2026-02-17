package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"subscription_service/pkg/logger"
)

func NewHandler(log logger.Logger, h *SubscriptionHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer, logger.GetLogMiddleware(log))

	r.Head("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/", h.ListSubscriptions)
		r.Get("/total", h.TotalSubscriptions)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSubscription)
			r.Put("/", h.UpdateSubscription)
			r.Delete("/", h.DeleteSubscription)
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
