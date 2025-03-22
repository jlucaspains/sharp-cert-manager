package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jlucaspains/sharp-cert-manager/handlers"
	"github.com/jlucaspains/sharp-cert-manager/jobs"
	"github.com/jlucaspains/sharp-cert-manager/midlewares"
	"github.com/jlucaspains/sharp-cert-manager/models"
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
	messageMentions, _ := os.LookupEnv("MESSAGE_MENTIONS")

	result.Init(jobs.Notifiers[webhookType], WebhookUrl, messageTitle, messageBody, messageUrl, messageMentions)

	return result
}

func startJobs(siteList []models.CheckCertItem) {
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

func getCORSOrigins() string {
	corsOrigins, ok := os.LookupEnv("CORS_ORIGINS")
	if ok {
		return corsOrigins
	}

	return ""
}

func stopJobs() {
	checkCertJob.Stop()
}

func startWebServer(siteList []models.CheckCertItem) {
	headless, ok := os.LookupEnv("HEADLESS")

	if !ok || headless == "true" {
		log.Println("Running in headless mode. Skipping web server start.")
		return
	}

	handlers := &handlers.Handlers{}
	handlers.CertList = siteList
	handlers.ExpirationWarningDays = getCertExpirationWarningDays()
	handlers.CORSOrigins = getCORSOrigins()

	router := http.NewServeMux()

	router.HandleFunc("GET /api/check-cert", handlers.CheckStatus)
	router.HandleFunc("GET /api/cert-list", handlers.GetCertList)
	router.HandleFunc("GET /health", handlers.HealthCheck)

	if handlers.CORSOrigins != "" {
		router.HandleFunc("OPTIONS /api/", handlers.CORS)
	}

	router.Handle("/", http.FileServer(http.Dir("./public/")))

	logRouter := midlewares.NewLogger(router)

	hostPort, ok := os.LookupEnv("WEB_HOST_PORT")
	if !ok {
		hostPort = ":8000"
	}

	certFile, useTls := os.LookupEnv("TLS_CERT_FILE")

	certKeyFile, ok := os.LookupEnv("TLS_CERT_KEY_FILE")
	useTls = useTls && ok

	log.Printf("Starting TLS server on port: %s; use tls: %t", hostPort, useTls)

	srv := &http.Server{
		Addr: hostPort,
	}

	srv.Handler = logRouter

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

func runOnce(siteList []models.CheckCertItem, done chan os.Signal) {
	schedule, _ := os.LookupEnv("CHECK_CERT_JOB_SCHEDULE")
	headless, _ := os.LookupEnv("HEADLESS")

	if schedule != "" || headless != "true" {
		return
	}

	log.Print("Running the checkCertJob once")
	level, _ := os.LookupEnv("CHECK_CERT_JOB_NOTIFICATION_LEVEL")
	warningDays := getCertExpirationWarningDays()
	err := checkCertJob.Init("* * * * *", level, warningDays, siteList, getJobNotifier())
	if err == nil {
		checkCertJob.RunNow()
	} else {
		log.Fatalf("Error running the checkCertJob once: %s", err)
	}

	// We don't want to wait on anything so do a graceful exist
	done <- syscall.SIGQUIT
}

func main() {
	loadEnv()

	siteList := shared.GetConfigCerts()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	startWebServer(siteList)
	startJobs(siteList)
	runOnce(siteList, done)

	<-done
	log.Print("Stopping jobs...")
	stopJobs()

	log.Print("All done. Bye!")
}
