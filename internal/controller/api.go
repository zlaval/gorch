package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorch/pkg/command"
	"gorch/pkg/rest"
	"log"
	"net/http"
)

type Api struct {
}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Run(ctx context.Context) error {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
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

	mux.Post("/command", a.command)

	return mux
}

func (a *Api) command(w http.ResponseWriter, r *http.Request) {
	var cmd command.Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(cmd)
	res := command.Response{Data: []string{"Test data"}}
	rest.SuccessResponse(w, res)
}
