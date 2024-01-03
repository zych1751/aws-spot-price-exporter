package main

import (
	"heartpattern.io/aws-spot-price-exporter/internal"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg := internal.ReadConfigOrDie()

	spotPriceExporter := internal.NewSpotPriceExporter(cfg)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(spotPriceExporter.RenderPrometheus(r.Context())))

		if err != nil {
			log.Printf("unable to write response, %v", err)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))

		if err != nil {
			log.Printf("unable to write response, %v", err)
		}
	})

	http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))

		if err != nil {
			log.Printf("unable to write response, %v", err)
		}
	})

	log.Printf("listening on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil))
}
