package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-manager/handlers"
	"github.com/jlucaspains/sharp-cert-manager/jobs"
	"github.com/jlucaspains/sharp-cert-manager/midlewares"
	"github.com/jlucaspains/sharp-cert-manager/shared"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var checkCertJob = &jobs.CheckCertJob{}
var env string

func loadEnv() {
	// outside of local environment, variables should be
	// OS environment variables
	env = os.Getenv("ENV")
	if err := godotenv.Load(); err != nil && env == "" {
		log.Fatal(fmt.Printf("Error loading .env file: %s", err))
	}
	env = os.Getenv("ENV") // reload env from .env file
}

func getJobNotifier() jobs.Notifier {
	result := &jobs.WebHookNotifier{}

	webhookType, _ := os.LookupEnv("WEBHOOK_TYPE")
	WebhookUrl, _ := os.LookupEnv("WEBHOOK_URL")
	messageUrl, _ := os.LookupEnv("MESSAGE_URL")
	messageTitle, _ := os.LookupEnv("MESSAGE_TITLE")
	messageBody, _ := os.LookupEnv("MESSAGE_BODY")

	result.Init(jobs.Notifiers[webhookType], WebhookUrl, messageTitle, messageBody, messageUrl)

	return result
}

func startJobs(siteList []string) {
	schedule, ok := os.LookupEnv("CHECK_CERT_JOB_SCHEDULE")

	if ok {
		log.Printf("Starting job engine with cron: %s", schedule)
		level, _ := os.LookupEnv("CHECK_CERT_JOB_NOTIFICATION_LEVEL")
		warningDays := getCertExpirationWarningDays()
		err := checkCertJob.Init(schedule, level, warningDays, siteList, getJobNotifier())
		if err == nil {
			checkCertJob.Start()
			log.Print("Job engine started")
		} else {
			log.Printf("Error starting job: %s", err)
		}
	} else {
		log.Println("No schedule defined for jobs")
	}
}

func getCertExpirationWarningDays() int {
	warningDaysConfig, _ := os.LookupEnv("CERT_WARNING_VALIDITY_DAYS")
	warningDays, _ := strconv.Atoi(warningDaysConfig)

	if warningDays > 0 {
		return warningDays
	}

	return 30
}

func stopJobs() {
	checkCertJob.Stop()
}

func startWebServer(siteList []string) {
	handlers := &handlers.Handlers{}
	handlers.SiteList = siteList
	handlers.ExpirationWarningDays = getCertExpirationWarningDays()

	router := mux.NewRouter()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, method string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		router.Handle(pattern, handler).Methods(method)
	}

	handleFunc("/api/check-url", "GET", handlers.CheckStatus)
	handleFunc("/api/site-list", "GET", handlers.GetSiteList)
	handleFunc("/health", "GET", handlers.HealthCheck)

	router.PathPrefix("/").Handler(otelhttp.NewHandler(http.FileServer(http.Dir("./public/")), "html"))

	logMiddleware := midlewares.NewLogMiddleware(log.Default())
	router.Use(logMiddleware.Func())

	hostPort, ok := os.LookupEnv("WEB_HOST_PORT")
	if !ok {
		hostPort = ":8000"
	}

	useTls := false
	certFile, ok := os.LookupEnv("TLS_CERT_FILE")
	useTls = ok

	certKeyFile, ok := os.LookupEnv("TLS_CERT_KEY_FILE")
	useTls = useTls && ok

	log.Printf("Starting TLS server on port: %s; use tls: %t", hostPort, useTls)

	srv := &http.Server{
		Addr: hostPort,
	}

	// if env == "local" {
	// 	srv.Handler = muxHandlers.CORS()(router)
	// } else {
	handler := otelhttp.NewHandler(router, "/")
	srv.Handler = handler
	// }

	go func() {
		var err error = nil
		if useTls {
			err = srv.ListenAndServeTLS(certFile, certKeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Web Server Started")
}

func main() {
	loadEnv()

	// Set up OpenTelemetry.
	serviceName := "sharp-cert-manager"
	serviceVersion := "1.0.0"
	otelShutdown, err := shared.SetupOTelSDK(context.Background(), serviceName, serviceVersion)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	siteList := shared.GetConfigSites()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startWebServer(siteList)
	startJobs(siteList)

	<-done
	log.Print("Stopping jobs...")
	stopJobs()

	log.Print("All done. Bye!")
}
