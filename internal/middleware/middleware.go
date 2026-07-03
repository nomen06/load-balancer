package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v", err)
				http.Error(w, "500-Internal server err", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
