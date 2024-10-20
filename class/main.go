package main

import (
	"class/config"
	"class/consumers"
	"class/models"
	"class/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// @title School Management API Leecho
// @version 0.1
// @description API for a school management system.
// @contact.name Akme
// @contact.url
// @contact.email
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitMQConfig, err := config.InitRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer rabbitMQConfig.Close()

	if err := rabbitMQConfig.DeclareQueue("class_created", true); err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	if err := db.AutoMigrate(&models.Class{}); err != nil {
		log.Fatalf("Failed to run migrations: %s", err)
	}

	consumers.StartClassCreatedConsumer(rabbitMQConfig, db)

	app := fiber.New()
	app.Static("/docs", "./public/")

	routes.ClassRoutes(app, rabbitMQConfig, db)

	log.Println("Starting server on :3000...")
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
