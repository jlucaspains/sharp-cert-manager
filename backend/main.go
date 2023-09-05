package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jlucaspains/sharp-cert-checker/handlers"
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

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	handlers := &handlers.Handlers{}
	handlers.SiteList = handlers.GetConfigSites()
	app.Get("/api/check-url", handlers.CheckStatus)
	app.Get("/api/site-list", handlers.GetSiteList)

	app.Static("/", "./public")

	app.Listen(":3000")
}
