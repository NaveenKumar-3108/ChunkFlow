package router

import (
	"CHUNKFLOW/handlers"

	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(RateLimiterMiddleware(rate.Limit(0.5), 1))

	r.HandleFunc("/upload", handlers.UploadAudio).Methods("POST")
	r.HandleFunc("/chunks{id}", handlers.GetChunkMetadata).Methods("GET")
	r.HandleFunc("/session{id}", handlers.GetUserChunksdata).Methods("GET")
	//	r.HandleFunc("/ws", handlers.HandleWebSocket)

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
