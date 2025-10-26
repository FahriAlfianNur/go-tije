package domain

import (
	"time"
)

type VehicleLocation struct {
	ID        string    `json:"id" db:"id"`
	VehicleID string    `json:"vehicle_id" db:"vehicle_id"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Timestamp int64     `json:"timestamp" db:"timestamp"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LocationMessage struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type GeofenceEvent struct {
	VehicleID string   `json:"vehicle_id"`
	Event     string   `json:"event"`
	Location  Location `json:"location"`
	Timestamp int64    `json:"timestamp"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Geofence struct {
	ID        string
	Name      string
	Latitude  float64
	Longitude float64
	Radius    float64
}