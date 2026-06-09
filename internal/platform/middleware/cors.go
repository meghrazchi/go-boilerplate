package middleware

import (
	"net/http"
	"strings"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
)

func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	allowedOrigins := make(map[string]struct{}, len(cfg.CORSAllowedOrigins))
	allowAll := false
	for _, origin := range cfg.CORSAllowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "*" {
			allowAll = true
		}
		if origin != "" {
			allowedOrigins[origin] = struct{}{}
		}
	}

	methods := strings.Join(cfg.CORSAllowedMethods, ", ")
	headers := strings.Join(cfg.CORSAllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if _, ok := allowedOrigins[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Add("Vary", "Origin")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
