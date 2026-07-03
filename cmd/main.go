package main

import (
	"log"
	"net/http"

	"github.com/nomen06/load-balancer/internal/proxy"
)

func main() {
	target := "http://localhost:8081" // will fix this hardcoding later
	p, err := proxy.NewProxy(target)
	if err != nil {
		log.Fatalf("initialisation failed for proxy : %v", err)
	}
	http.Handle("/", p)
	log.Println("Reverse proxy initialisation on our port:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
