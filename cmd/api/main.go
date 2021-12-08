package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const vesion = "1.0.0"

type config struct {
	port int
	env  string
}

type Appstatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"vesion"`
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.Parse()

	fmt.Println("Running ")

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		currentStatus := Appstatus{
			Status:      "Available",
			Environment: cfg.env,
			Version:     vesion,
		}

		js, err := json.MarshalIndent(currentStatus, "", "\t")

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(js)

	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)

	if err != nil {
		log.Println(err)
	}

}
