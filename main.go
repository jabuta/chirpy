package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	cfgapi := &apiConfig{fileserverHits: 0}

	mainR := chi.NewRouter()

	fsHandler := cfgapi.middleWareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mainR.Handle("/app", fsHandler)
	mainR.Handle("/app/*", fsHandler)

	apiR := chi.NewRouter()
	apiR.Get("/healthz", healthzHandler)

	apiR.Get("/reset", cfgapi.metricsResetHandler)
	mainR.Mount("/api", apiR)

	adminR := chi.NewRouter()
	adminR.Get("/metrics", cfgapi.metricsHandler)
	mainR.Mount("/admin", adminR)

	corsr := middlewareCors(mainR)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsr,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
