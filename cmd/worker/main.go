package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fahri/go-tije/internal/config"
	"github.com/fahri/go-tije/internal/domain"
	"github.com/streadway/amqp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		cfg.RabbitMQ.Exchange, 
		"topic",               
		true,                  
		false,                 
		false,                 
		false,                 
		nil,                   
	)
	if err != nil {
		log.Fatal("Failed to declare an exchange:", err)
	}

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queue, 
		true,               
		false,              
		false,              
		false,              
		nil,                
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	err = ch.QueueBind(
		q.Name,               
		"geofence.#",          
		cfg.RabbitMQ.Exchange, 
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to bind queue to exchange:", err)
	}

	err = ch.Qos(
		1,     
		0,     
		false, 
	)
	if err != nil {
		log.Fatal("Failed to set QoS:", err)
	}

	msgs, err := ch.Consume(
		q.Name,       
		"worker",     
		false,         
		false,         
		false,         
		false,         
		nil,           
	)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Worker started. Waiting for geofence events...")
	log.Printf("Connected to queue: %s", cfg.RabbitMQ.Queue)

	go func() {
		for msg := range msgs {
			processMessage(msg)
		}
	}()

	<-sigChan
	log.Println("Shutting down worker...")
}

func processMessage(msg amqp.Delivery) {
	var event domain.GeofenceEvent
	
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Failed to parse message: %v", err)
		msg.Reject(false)
		return
	}

	log.Printf("=== GEOFENCE ALERT ===")
	log.Printf("Vehicle ID: %s", event.VehicleID)
	log.Printf("Event Type: %s", event.Event)
	log.Printf("Location: %.4f, %.4f", event.Location.Latitude, event.Location.Longitude)
	log.Printf("Timestamp: %s", time.Unix(event.Timestamp, 0).Format("2006-01-02 15:04:05"))
	log.Printf("=====================")

	
	time.Sleep(100 * time.Millisecond)

	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	} else {
		log.Printf("Message processed successfully for vehicle %s", event.VehicleID)
	}
}