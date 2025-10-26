package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fahri/go-tije/internal/config"
	"github.com/fahri/go-tije/internal/repository"
	"github.com/fahri/go-tije/internal/service"
	mqttclient "github.com/fahri/go-tije/pkg/mqtt"
	"github.com/fahri/go-tije/pkg/rabbitmq"
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
	
	mqttClient, err := mqttclient.NewClient(&cfg.MQTT)
	if err != nil {
		log.Fatal("Failed to connect to MQTT:", err)
	}
	defer mqttClient.Disconnect()
	
	topic := "/fleet/vehicle/+/location"
	handler := func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received location data from topic %s", msg.Topic())
		
		if err := vehicleService.ProcessLocation(context.Background(), msg.Payload()); err != nil {
			log.Printf("Failed to process location: %v", err)
		} else {
			log.Printf("Location processed successfully")
		}
	}
	
	if err := mqttClient.Subscribe(topic, handler); err != nil {
		log.Fatal("Failed to subscribe to topic:", err)
	}
	
	fmt.Printf("MQTT Subscriber started. Listening to topic: %s\n", topic)
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	log.Println("Shutting down...")
}