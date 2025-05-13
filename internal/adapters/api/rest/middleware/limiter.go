package middleware

import (
	"balancer/internal/adapters/api/rest/erros"
	"balancer/internal/logger"
	"balancer/internal/ratelimit"
	"net/http"
)

func Middleware(limiter *ratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := r.Header.Get("X-Client-ID")
			if clientID == "" {
				clientID = r.RemoteAddr
			}

			if !limiter.Allow(clientID) {
				logger.Warn(r.Context(), "rate limit exceeded", "clientID", clientID)
				erros.JSON(w, http.StatusTooManyRequests, "Rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
