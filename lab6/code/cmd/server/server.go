package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("WOW, I got a request on `\\health` endpoint !!!")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Hello, I'm healthy"))
		if err != nil {
			log.Printf("Failed to write health check response: %v\n", err)
		}
	})

	log.Printf("Server started on port: %d", *port)
	err = http.Serve(l, nil)
	if err != nil {
		panic(err)
	}
}
