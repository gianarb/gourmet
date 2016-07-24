package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type Ping struct {
	Status bool              `json:"status"`
	Info   map[string]string `json:"info"`
}

func HealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("New ping")
		ping := Ping{
			Status: true,
		}

		checkDockerClient(&ping)

		js, _ := json.Marshal(ping)
		w.Header().Set("Content-Type", "application/json")
		if len(ping.Info) > 0 {
			ping.Status = false
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
		w.Write(js)
	}
}

func checkDockerClient(ping *Ping) {
	_, err := docker.NewClientFromEnv()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": "ping",
			"error":   err,
		}).Warn("Problem to communicate with docker")
		ping.Info["docker"] = err.Error()
	}
}
