package router

import (
	"CHUNKFLOW/handlers"
	"log"
	"runtime/debug"

	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(RecoverMiddleware)
	r.Use(RateLimiterMiddleware(rate.Limit(0.5), 1))

	r.HandleFunc("/upload", handlers.UploadAudio).Methods("POST")
	r.HandleFunc("/chunks{id}", handlers.GetChunkMetadata).Methods("GET")
	r.HandleFunc("/session{id}", handlers.GetUserChunksdata).Methods("GET")
	r.HandleFunc("/ws", handlers.AudioWebSocket)

	return r
}

func RateLimiterMiddleware(limit rate.Limit, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(limit, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				log.Printf("Panic recovered: %v\nStack trace:\n%s", err, debug.Stack())

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"status":"E","errmsg":"Internal Server Error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
