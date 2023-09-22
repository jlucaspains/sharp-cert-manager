package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	muxHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jlucaspains/sharp-cert-manager/handlers"
	"github.com/jlucaspains/sharp-cert-manager/jobs"
	"github.com/jlucaspains/sharp-cert-manager/midlewares"
	"github.com/jlucaspains/sharp-cert-manager/shared"
	"github.com/joho/godotenv"
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
	router.HandleFunc("/api/check-url", handlers.CheckStatus).Methods("GET")
	router.HandleFunc("/api/site-list", handlers.GetSiteList).Methods("GET")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

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

	if env == "local" {
		srv.Handler = muxHandlers.CORS()(router)
	} else {
		srv.Handler = router
	}

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
