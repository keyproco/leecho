package main

import (
	"leecho/config"
	"leecho/routes"
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

	rabbitMQConfig := config.InitRabbitMQ()
	defer rabbitMQConfig.Close()

	config.ConnectDatabase()
	app := fiber.New()
	app.Static("/docs", "./public/")

	routes.ClassRoutes(app)

	app.Listen(":3000")
}
