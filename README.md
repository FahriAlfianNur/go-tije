# Fleet Management System

Real-time vehicle tracking system with geofence detection capabilities.

## Features

- Real-time vehicle location tracking via MQTT
- Location history storage in PostgreSQL
- REST API for data retrieval
- Automatic geofence detection and alerting
- Containerized deployment with Docker

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Setup with Docker

1. Clone the repository
```bash
git clone https://github.com/fahri/go-tije.git
cd go-tije
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Start services
```bash
docker-compose up -d
```

4. Verify services are running
```bash
docker-compose ps
```

### Local Development

1. Install dependencies
```bash
go mod download
```

2. Set up environment
```bash
cp .env.example .env
```

3. Start PostgreSQL, RabbitMQ, and Mosquitto
```bash
docker-compose up -d postgres rabbitmq mosquitto
```

4. Run database migrations
```bash
psql -h localhost -U fleet_user -d fleet_db < migrations/init.sql
```

5. Start services
```bash
# Terminal 1: API Server
go run cmd/api/main.go

# Terminal 2: MQTT Subscriber
go run cmd/subscriber/main.go

# Terminal 3: Vehicle Simulator
go run cmd/publisher/main.go
```

## API Endpoints

### Get Latest Vehicle Location
```bash
GET /vehicles/{vehicle_id}/location

curl http://localhost:8080/vehicles/B1234XYZ/location
```

Response:
```json
{
  "id": "uuid",
  "vehicle_id": "B1234XYZ",
  "latitude": -6.2088,
  "longitude": 106.8456,
  "timestamp": 1715003456,
  "created_at": "2024-05-06T12:00:00Z"
}
```

### Get Location History
```bash
GET /vehicles/{vehicle_id}/history?start={timestamp}&end={timestamp}

curl "http://localhost:8080/vehicles/B1234XYZ/history?start=1715000000&end=1715009999"
```

Response:
```json
[
  {
    "id": "uuid",
    "vehicle_id": "B1234XYZ",
    "latitude": -6.2088,
    "longitude": 106.8456,
    "timestamp": 1715003456,
    "created_at": "2024-05-06T12:00:00Z"
  }
]
```

## Testing

### Manual Testing

#### 1. Verify Data Flow

**A. Check Publisher Sending Data**
```bash
docker logs fleet-publisher --tail=5
```

**B. Check Subscriber Receiving and Storing**
```bash
docker logs fleet-subscriber --tail=5
```

**C. Verify Data in Database**
```bash
docker exec fleet-postgres psql -U fleet_user -d fleet_db -c "SELECT COUNT(*) FROM vehicle_locations;"
```

#### 2. Test API Endpoints

**A. Latest Vehicle Location**
```bash
curl http://localhost:8080/vehicles/B1234XYZ/location | jq .
```

**B. Location History (Time Range)**
```bash
NOW=$(date +%s)
PAST=$((NOW - 600))
curl "http://localhost:8080/vehicles/B1234XYZ/history?start=$PAST&end=$NOW" | jq '. | length'
```

#### 3. Verify RabbitMQ & Worker

**A. Check Exchange and Queue**
```bash
docker exec fleet-rabbitmq rabbitmqctl list_exchanges | grep fleet
```

**B. Check Queue and Consumer**
```bash
docker exec fleet-rabbitmq rabbitmqctl list_queues name consumers
```

**C. Check Worker Logs**
```bash
docker logs fleet-worker --tail=5
```


## Service Ports

| Service     | Port  | Description           |
|------------|-------|----------------------|
| API        | 8080  | REST API             |
| PostgreSQL | 5432  | Database             |
| RabbitMQ   | 5672  | AMQP                 |
| RabbitMQ   | 15672 | Management UI        |
| Mosquitto  | 1883  | MQTT                 |
| Mosquitto  | 9001  | WebSocket            |

## Configuration

Environment variables can be configured in `.env` file:

- `APP_PORT`: API server port (default: 8080)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `MQTT_BROKER`: MQTT broker URL
- `RABBITMQ_URL`: RabbitMQ connection URL
- `GEOFENCE_RADIUS`: Detection radius in meters
- `GEOFENCE_LAT`: Geofence center latitude
- `GEOFENCE_LON`: Geofence center longitude



## License

MIT