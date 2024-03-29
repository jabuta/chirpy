package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jabuta/chirpy/internal/database"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
	polka_key      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	const defaultDbPath = "database.json"
	godotenv.Load()

	dbPath := debugMode(defaultDbPath)

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %s", err)
	}
	cfgapi := &apiConfig{
		fileserverHits: 0,
		db:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polka_key:      os.Getenv("POLKA_API"),
	}

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
	apiR.Get("/chirps", cfgapi.getChirps)
	apiR.Post("/chirps", cfgapi.postChirp)
	apiR.Get("/chirps/{chirpID}", cfgapi.getChirp)
	apiR.Delete("/chirps/{chirpID}", cfgapi.deleteChirp)

	//api_users.go
	apiR.Post("/users", cfgapi.newUser)
	apiR.Put("/users", cfgapi.putUser)
	apiR.Post("/login", cfgapi.loginUser)
	apiR.Post("/refresh", cfgapi.apiRefresh)
	apiR.Post("/revoke", cfgapi.apiRevoke)

	//api_polka.go
	apiR.Post("/polka/webhooks", cfgapi.polkaWebhooks)

	mainR.Mount("/api", apiR)

	//metrics router metrics.goS
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
