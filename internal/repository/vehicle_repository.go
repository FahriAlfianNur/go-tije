package repository

import (
	"context"
	"fmt"

	"github.com/fahri/go-tije/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleRepository interface {
	Save(ctx context.Context, location *domain.VehicleLocation) error
	FindLatest(ctx context.Context, vehicleID string) (*domain.VehicleLocation, error)
	FindHistory(ctx context.Context, vehicleID string, start, end int64) ([]*domain.VehicleLocation, error)
}

type vehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) Save(ctx context.Context, location *domain.VehicleLocation) error {
	query := `
		INSERT INTO vehicle_locations (id, vehicle_id, latitude, longitude, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`
	
	location.ID = uuid.New().String()
	_, err := r.db.Exec(ctx, query,
		location.ID,
		location.VehicleID,
		location.Latitude,
		location.Longitude,
		location.Timestamp,
	)
	
	return err
}

func (r *vehicleRepository) FindLatest(ctx context.Context, vehicleID string) (*domain.VehicleLocation, error) {
	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`
	
	var location domain.VehicleLocation
	err := r.db.QueryRow(ctx, query, vehicleID).Scan(
		&location.ID,
		&location.VehicleID,
		&location.Latitude,
		&location.Longitude,
		&location.Timestamp,
		&location.CreatedAt,
	)
	
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vehicle not found")
	}
	
	return &location, err
}

func (r *vehicleRepository) FindHistory(ctx context.Context, vehicleID string, start, end int64) ([]*domain.VehicleLocation, error) {
	query := `
		SELECT id, vehicle_id, latitude, longitude, timestamp, created_at
		FROM vehicle_locations
		WHERE vehicle_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC
	`
	
	rows, err := r.db.Query(ctx, query, vehicleID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var locations []*domain.VehicleLocation
	for rows.Next() {
		var location domain.VehicleLocation
		err := rows.Scan(
			&location.ID,
			&location.VehicleID,
			&location.Latitude,
			&location.Longitude,
			&location.Timestamp,
			&location.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}
	
	return locations, nil
}