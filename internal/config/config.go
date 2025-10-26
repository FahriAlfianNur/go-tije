package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	DB        DatabaseConfig
	MQTT      MQTTConfig
	RabbitMQ  RabbitMQConfig
	Geofence  GeofenceConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type MQTTConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

type RabbitMQConfig struct {
	URL      string
	Exchange string
	Queue    string
}

type GeofenceConfig struct {
	Radius    float64
	Latitude  float64
	Longitude float64
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// Continue without .env file
	}

	geofenceRadius, _ := strconv.ParseFloat(getEnv("GEOFENCE_RADIUS", "50"), 64)
	geofenceLat, _ := strconv.ParseFloat(getEnv("GEOFENCE_LAT", "-6.2088"), 64)
	geofenceLon, _ := strconv.ParseFloat(getEnv("GEOFENCE_LON", "106.8456"), 64)

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "fleet_user"),
			Password: getEnv("DB_PASSWORD", "fleet_password"),
			Name:     getEnv("DB_NAME", "fleet_db"),
		},
		MQTT: MQTTConfig{
			Broker:   getEnv("MQTT_BROKER", "tcp://localhost:1883"),
			ClientID: getEnv("MQTT_CLIENT_ID", "fleet-subscriber"),
			Username: getEnv("MQTT_USERNAME", ""),
			Password: getEnv("MQTT_PASSWORD", ""),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			Exchange: getEnv("RABBITMQ_EXCHANGE", "fleet.events"),
			Queue:    getEnv("RABBITMQ_QUEUE", "geofence_alerts"),
		},
		Geofence: GeofenceConfig{
			Radius:    geofenceRadius,
			Latitude:  geofenceLat,
			Longitude: geofenceLon,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}