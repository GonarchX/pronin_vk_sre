package main

import (
	"log"
	"math/rand/v2"
	"net/http"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if rand.IntN(2) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("BAD"))
			if err != nil {
				panic(err)
			}
			log.Println("Get request with '/status' path\nRespond BAD")
		} else {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				panic(err)
			}
			log.Println("Get request with '/status' path\nRespond OK")
		}
	})

	log.Println("Start server on port :8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
