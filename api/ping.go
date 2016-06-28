package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type Ping struct {
	Status string
}

func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("New ping")
		ping := Ping{"ok"}
		js, _ := json.Marshal(ping)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(js)
	}
}
