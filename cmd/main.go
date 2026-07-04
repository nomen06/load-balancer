package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nomen06/load-balancer/internal/balancer"
	"github.com/nomen06/load-balancer/internal/config"
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
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	serverpool := &balancer.Serverpool{}
	// servers := []string{
	// 	"http://localhost:8081",
	// 	"http://localhost:8082",
	// 	"http://localhost:8083",
	// }
	for _, srv := range cfg.Servers {
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

	// log.Println("Load balancer running on : 8080")
	go func() {
		t := time.NewTicker(10 * time.Second)
		for range t.C {
			log.Println("Running health check...")
			serverpool.Healthcheck()
		}
	}()
	serverAddr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}
	go func() {
		log.Printf("Load balaner running on %s", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
	// log.Printf("Load balancer running on %s", serverAddr)

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatalf("Server failed: %v", err)
	// }

	// channel listening to signal terminations
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Termination signal recieved. Server shutting down gracefully..")

	// giving a deadline context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly.")
}
