package main

import (
	"log"

	"github.com/fahri/go-tije/internal/config"
	"github.com/fahri/go-tije/internal/handler"
	"github.com/fahri/go-tije/internal/repository"
	"github.com/fahri/go-tije/internal/service"
	"github.com/fahri/go-tije/pkg/rabbitmq"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	
	db, err := repository.NewDatabase(&cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	rmqPublisher, err := rabbitmq.NewPublisher(&cfg.RabbitMQ)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rmqPublisher.Close()
	
	vehicleRepo := repository.NewVehicleRepository(db)
	vehicleService := service.NewVehicleService(vehicleRepo, rmqPublisher, &cfg.Geofence)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	
	app := fiber.New()
	
	app.Use(logger.New())
	app.Use(cors.New())
	
	api := app.Group("/vehicles")
	api.Get("/:vehicle_id/location", vehicleHandler.GetLatestLocation)
	api.Get("/:vehicle_id/history", vehicleHandler.GetLocationHistory)
	
	log.Printf("Server starting on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}