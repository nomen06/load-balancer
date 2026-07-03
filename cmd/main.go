package main

import (
	"log"
	"net/http"
	"time"

	"github.com/nomen06/load-balancer/internal/balancer"
	"github.com/nomen06/load-balancer/internal/middleware"
)

func main() {
	// target := "http://localhost:8081" // will fix this hardcoding later
	// p, err := proxy.NewProxy(target)
	// if err != nil {
	// 	log.Fatalf("initialisation failed for proxy : %v", err)
	// }
	// http.Handle("/", p)
	// log.Println("Reverse proxy initialisation on our port:8080")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatalf("server failed: %v", err)
	// }

	//connecting for multiple backends
	serverpool := &balancer.Serverpool{}
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	for _, srv := range servers {
		if err := serverpool.AddBackend(srv); err != nil {
			log.Fatalf("Could not parse backend URL : %v", err)
		}
		log.Printf("Configured backend: %s", srv)
	}
	// http.Handle("/", serverpool) instead of this, now we'lll wrap our pool in middleware chain of recovery,logging,serverpool
	var handler http.Handler = serverpool
	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)
	http.Handle("/", handler) //yaayyyy testing time now

	log.Println("Load balancer running on : 8080")
	go func() {
		t := time.NewTicker(10 * time.Second)
		for range t.C {
			log.Println("Running health check...")
			serverpool.Healthcheck()
		}
	}()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
