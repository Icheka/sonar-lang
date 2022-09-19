package server

import (
	"fmt"
	"language-server/server/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func Start() {
	router := mux.NewRouter()

	// register middleware
	router.Use(enableCORS)

	// register handlers
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi!"))
	})
	router.HandleFunc("/ws", handlers.HandleWebSocket)

	port := 9999
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
