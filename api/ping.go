package api

import (
    "net/http"
	"encoding/json"
	"log"
)

type Ping struct {
	Status string
}

func PingHandler (logger *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return  func(w http.ResponseWriter, r *http.Request) {
		ping := Ping{"ok"}
		js, _ := json.Marshal(ping)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
		logger.Printf("Ping")
	}
}

