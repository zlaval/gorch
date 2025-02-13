package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorch/pkg/command"
	"gorch/pkg/registry"
	"gorch/pkg/rest"
	"log"
	"log/slog"
	"net/http"
)

type Api struct {
	db *Database
	sc *Scheduler
}

func NewApi(db *Database, sc *Scheduler) *Api {
	return &Api{
		db: db,
		sc: sc,
	}
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

	slog.Info("Server is up and running")
	<-ctx.Done()

	return s.Shutdown(ctx)
}

func (a *Api) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Post("/command", a.command)
	mux.Post("/register", a.registerWorker)

	return mux
}

func (a *Api) registerWorker(w http.ResponseWriter, r *http.Request) {
	var wr registry.WorkerRegRequest
	if err := json.NewDecoder(r.Body).Decode(&wr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	worker := WorkerEntity{
		Name:    wr.Name,
		Address: wr.Address,
		IP:      wr.IP,
		Port:    wr.Port,
		Status:  WorkerUp,
	}

	if err := a.db.SaveWorker(worker); err != nil {
		slog.Error("register worker", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("Worker has been registered", slog.String("name", wr.Name))

	rest.SuccessResponse(w, "OK")
}

func (a *Api) command(w http.ResponseWriter, r *http.Request) {
	var cmd command.Command
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO remove this mock task
	fmt.Println(cmd.String())
	a.sc.Add(Task{})

	res := command.Response{Data: []string{"Test data"}}
	rest.SuccessResponse(w, res)
}
