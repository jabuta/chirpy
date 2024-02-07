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

	// main app router
	mainR := chi.NewRouter()
	fsHandler := cfgapi.middleWareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mainR.Handle("/app", fsHandler)
	mainR.Handle("/app/*", fsHandler)

	//api router api.go
	apiR := chi.NewRouter()
	apiR.Get("/healthz", healthzHandler)
	apiR.Get("/reset", cfgapi.metricsResetHandler)
	//api_chirps.go
	apiR.Post("/chirps", postChirpHandler)

	mainR.Mount("/api", apiR)

	//metrics router metrics.go
	adminR := chi.NewRouter()
	adminR.Get("/metrics", cfgapi.metricsHandler)
	mainR.Mount("/admin", adminR)

	//set headers
	corsr := middlewareCors(mainR)

	//start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsr,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
