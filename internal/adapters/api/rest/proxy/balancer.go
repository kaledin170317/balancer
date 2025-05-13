package proxy

import (
	"balancer/internal/adapters/api/rest/erros"
	"balancer/internal/balancer"
	"balancer/internal/logger"
	"net/http"
	"net/http/httputil"
)

func NewBalancer(balancer balancer.BalanceAlgorithm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		backend := balancer.Next()
		if backend == nil {
			log.Warn("no healthy backend available")
			erros.JSON(w, http.StatusServiceUnavailable, "No healthy backend")
			return
		}

		backend.IncConnections()
		defer backend.DecConnections()

		log.Info("proxying request", "target", backend.URL.String())

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Warn("backend failed", "backend", backend.URL.String(), "err", err)
			backend.SetAlive(false)
			erros.JSON(rw, http.StatusBadGateway, "Backend error")
		}

		proxy.ServeHTTP(w, r)
	}
}
