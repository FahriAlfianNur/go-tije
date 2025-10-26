package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/fahri/go-tije/internal/config"
	"github.com/fahri/go-tije/internal/domain"
	"github.com/fahri/go-tije/internal/repository"
	"github.com/fahri/go-tije/pkg/geofence"
	"github.com/fahri/go-tije/pkg/rabbitmq"
)

type VehicleService interface {
	ProcessLocation(ctx context.Context, message []byte) error
	GetLatestLocation(ctx context.Context, vehicleID string) (*domain.VehicleLocation, error)
	GetLocationHistory(ctx context.Context, vehicleID string, start, end int64) ([]*domain.VehicleLocation, error)
}

type vehicleService struct {
	repo      repository.VehicleRepository
	rabbitmq  *rabbitmq.Publisher
	geofenceConfig *config.GeofenceConfig
}

func NewVehicleService(repo repository.VehicleRepository, rmq *rabbitmq.Publisher, geofenceCfg *config.GeofenceConfig) VehicleService {
	return &vehicleService{
		repo:      repo,
		rabbitmq:  rmq,
		geofenceConfig: geofenceCfg,
	}
}

func (s *vehicleService) ProcessLocation(ctx context.Context, message []byte) error {
	var locationMsg domain.LocationMessage
	if err := json.Unmarshal(message, &locationMsg); err != nil {
		return err
	}
	
	location := &domain.VehicleLocation{
		VehicleID: locationMsg.VehicleID,
		Latitude:  locationMsg.Latitude,
		Longitude: locationMsg.Longitude,
		Timestamp: locationMsg.Timestamp,
	}
	
	if err := s.repo.Save(ctx, location); err != nil {
		return err
	}
	
	if s.checkGeofence(location) {
		event := domain.GeofenceEvent{
			VehicleID: location.VehicleID,
			Event:     "geofence_entry",
			Location: domain.Location{
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
			Timestamp: location.Timestamp,
		}
		
		eventData, _ := json.Marshal(event)
		if err := s.rabbitmq.Publish(ctx, eventData); err != nil {
			log.Printf("Failed to publish geofence event: %v", err)
		}
	}
	
	return nil
}

func (s *vehicleService) GetLatestLocation(ctx context.Context, vehicleID string) (*domain.VehicleLocation, error) {
	return s.repo.FindLatest(ctx, vehicleID)
}

func (s *vehicleService) GetLocationHistory(ctx context.Context, vehicleID string, start, end int64) ([]*domain.VehicleLocation, error) {
	return s.repo.FindHistory(ctx, vehicleID, start, end)
}

func (s *vehicleService) checkGeofence(location *domain.VehicleLocation) bool {
	center := geofence.Point{
		Latitude:  s.geofenceConfig.Latitude,
		Longitude: s.geofenceConfig.Longitude,
	}
	
	target := geofence.Point{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
	
	return geofence.IsWithinRadius(center, target, s.geofenceConfig.Radius)
}