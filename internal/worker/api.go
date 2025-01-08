package worker

import (
	"Gorch/pkg/rest"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Api struct {
	worker Worker
}

func NewApi(worker Worker) *Api {
	return &Api{
		worker: worker,
	}
}

func (a *Api) Run(ctx context.Context) error {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8081),
		Handler: a.routes(),
	}

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Http server error: %v", err)
		}
	}()

	<-ctx.Done()

	return s.Shutdown(ctx)
}

func (a *Api) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/health", a.health)

	return mux
}

func (a *Api) health(w http.ResponseWriter, _ *http.Request) {
	rest.SuccessResponse(w, struct{ Status string }{"UP"})
}
