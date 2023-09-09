package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-checker/handlers"
	"github.com/jlucaspains/sharp-cert-checker/midlewares"
	"github.com/joho/godotenv"
)

func loadEnv() {
	// outside of local environment, variables should be
	// OS environment variables
	env := os.Getenv("ENV")
	if err := godotenv.Load(); err != nil && env == "" {
		log.Fatal(fmt.Printf("Error loading .env file: %s", err))
	}
}

func main() {
	loadEnv()

	handlers := &handlers.Handlers{}
	handlers.SiteList = handlers.GetConfigSites()

	router := mux.NewRouter()
	router.HandleFunc("/api/check-url", handlers.CheckStatus).Methods("GET")
	router.HandleFunc("/api/site-list", handlers.GetSiteList).Methods("GET")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	logMiddleware := midlewares.NewLogMiddleware(log.Default())
	router.Use(logMiddleware.Func())

	host_port, ok := os.LookupEnv("WEB_HOST_PORT")
	if !ok {
		host_port = ":3000"
	}

	use_tls := false
	cert_file, ok := os.LookupEnv("TLS_CERT_FILE")
	if ok {
		use_tls = true
	}

	cert_key_file, ok := os.LookupEnv("TLS_CERT_KEY_FILE")
	if ok {
		use_tls = use_tls && true
	}

	log.Printf("Starting TLS server on port: %s; use tls: %t", host_port, use_tls)
	if use_tls {
		log.Fatalln(http.ListenAndServeTLS(host_port, cert_file, cert_key_file, router))
	} else {
		log.Fatalln(http.ListenAndServe(host_port, router))
	}
}
