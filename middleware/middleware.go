package middleware

import (
	"net/http"
	"time"

	"github.com/jakabrajadenta/go-explore-api/pkg/logger"
)

type Func func(http.Handler) http.Handler

// Chain wraps a handler with the given middlewares in order (first = outermost).
func Chain(h http.Handler, middlewares ...Func) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// Logger generates a trace ID per request, injects it into the context, exposes
// it via the X-Trace-Id response header, and logs the request lifecycle.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := logger.Generate()
		ctx := logger.WithTraceID(r.Context(), traceID)
		r = r.WithContext(ctx)

		rw := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		rw.Header().Set("X-Trace-Id", traceID)

		start := time.Now()
		log := logger.FromCtx(ctx)
		log.Info("request received",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)

		next.ServeHTTP(rw, r)

		log.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"latency", time.Since(start).String(),
		)
	})
}

// CORS adds permissive CORS headers suitable for local development.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
