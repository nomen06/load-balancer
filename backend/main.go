package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "backend is running, port : %s\n", port)
		fmt.Printf("request recievd on %s\n", r.URL.Path)
	})
	log.Printf("backend server listening, port:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
