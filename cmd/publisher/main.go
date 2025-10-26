package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fahri/go-tije/internal/config"
	"github.com/fahri/go-tije/internal/domain"
	mqttclient "github.com/fahri/go-tije/pkg/mqtt"
)

type VehicleSimulator struct {
	vehicleID string
	lat       float64
	lon       float64
}

func NewVehicleSimulator(vehicleID string, startLat, startLon float64) *VehicleSimulator {
	return &VehicleSimulator{
		vehicleID: vehicleID,
		lat:       startLat,
		lon:       startLon,
	}
}

func (v *VehicleSimulator) move() {
	deltaLat := (rand.Float64() - 0.5) * 0.001
	deltaLon := (rand.Float64() - 0.5) * 0.001
	
	v.lat += deltaLat
	v.lon += deltaLon
	
	v.lat = math.Max(-90, math.Min(90, v.lat))
	v.lon = math.Max(-180, math.Min(180, v.lon))
}

func (v *VehicleSimulator) getLocation() domain.LocationMessage {
	return domain.LocationMessage{
		VehicleID: v.vehicleID,
		Latitude:  v.lat,
		Longitude: v.lon,
		Timestamp: time.Now().Unix(),
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	
	mqttClient, err := mqttclient.NewClient(&cfg.MQTT)
	if err != nil {
		log.Fatal("Failed to connect to MQTT:", err)
	}
	defer mqttClient.Disconnect()
	
	vehicles := []*VehicleSimulator{
		NewVehicleSimulator("B1234XYZ", -6.2088, 106.8456),
		NewVehicleSimulator("B5678ABC", -6.2100, 106.8470),
		NewVehicleSimulator("B9012DEF", -6.2050, 106.8430),
	}
	
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("Starting MQTT publisher...")
	
	for {
		select {
		case <-ticker.C:
			for _, vehicle := range vehicles {
				vehicle.move()
				location := vehicle.getLocation()
				
				data, err := json.Marshal(location)
				if err != nil {
					log.Printf("Failed to marshal location: %v", err)
					continue
				}
				
				topic := fmt.Sprintf("/fleet/vehicle/%s/location", vehicle.vehicleID)
				if err := mqttClient.Publish(topic, data); err != nil {
					log.Printf("Failed to publish to %s: %v", topic, err)
				} else {
					log.Printf("Published location for %s: lat=%.4f, lon=%.4f",
						vehicle.vehicleID, location.Latitude, location.Longitude)
				}
			}
			
		case <-sigChan:
			log.Println("Shutting down publisher...")
			return
		}
	}
}