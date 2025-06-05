package main

import (
	"net/http"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
    mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	apiCfg := &apiConfig{}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", apiCfg.getFileserverHits)
	mux.HandleFunc("/reset", apiCfg.resetMetrics)
    http.ListenAndServe(":8080", mux)

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFileserverHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load()),
	))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Metrics reset\n"))
	}