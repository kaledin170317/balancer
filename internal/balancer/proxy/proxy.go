package proxy

import (
	"balancer/internal/balancer"
	"balancer/internal/logger"
	"net/http"
	"net/http/httputil"
)

func NewHandler(balancer balancer.BalanceAlgorithm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		backend := balancer.Next()
		if backend == nil {
			log.Warn("no healthy backend available")
			http.Error(w, "no healthy backend", http.StatusServiceUnavailable)
			return
		}

		backend.IncConnections()
		defer backend.DecConnections()

		log.Info("proxying request", "target", backend.URL.String())

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Warn("backend failed", "backend", backend.URL.String(), "err", err)
			backend.SetAlive(false)
			http.Error(rw, "backend error", http.StatusBadGateway)
		}

		proxy.ServeHTTP(w, r)
	}
}
