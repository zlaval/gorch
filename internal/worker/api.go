package worker

import (
	"Gorch/pkg/pod"
	"Gorch/pkg/rest"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Api struct {
	worker           Worker
	metricsCollector MetricsCollector
}

func NewApi(worker Worker, metricsCollector MetricsCollector) *Api {
	return &Api{
		worker:           worker,
		metricsCollector: metricsCollector,
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
	mux.Get("/metrics", a.metrics)

	mux.Route("/pods", func(r chi.Router) {
		r.Post("/", a.createPod)
		r.Get("/", a.listPods)
		r.Get("/metrics", a.podMetrics)
		r.Route("/{containerID}", func(r chi.Router) {
			r.Delete("/", a.deletePod)
		})
	})

	return mux
}

func (a *Api) health(w http.ResponseWriter, _ *http.Request) {
	rest.SuccessResponse(w, struct{ Status string }{"UP"})
}

func (a *Api) metrics(w http.ResponseWriter, _ *http.Request) {
	res := a.metricsCollector.CollectNodeMetrics()
	rest.SuccessResponse(w, res)
}

func (a *Api) podMetrics(w http.ResponseWriter, r *http.Request) {
	res := a.metricsCollector.CollectPodMetrics(r.Context())
	rest.SuccessResponse(w, res)
}

func (a *Api) listPods(w http.ResponseWriter, r *http.Request) {
	res, err := a.worker.ListPods(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rest.SuccessResponse(w, res)
}

func (a *Api) createPod(w http.ResponseWriter, r *http.Request) {
	var cr pod.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := a.worker.Start(r.Context(), cr)
	rest.SuccessResponse(w, res)
}

func (a *Api) deletePod(w http.ResponseWriter, r *http.Request) {
	containerID := chi.URLParam(r, "containerID")
	res := a.worker.Stop(r.Context(), containerID)
	rest.SuccessResponse(w, res)
}
