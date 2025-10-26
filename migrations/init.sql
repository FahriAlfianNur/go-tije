CREATE TABLE IF NOT EXISTS vehicle_locations (
    id VARCHAR(36) PRIMARY KEY,
    vehicle_id VARCHAR(50) NOT NULL,
    latitude DECIMAL(10, 6) NOT NULL,
    longitude DECIMAL(10, 6) NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vehicle_id ON vehicle_locations(vehicle_id);
CREATE INDEX idx_timestamp ON vehicle_locations(timestamp);
CREATE INDEX idx_vehicle_timestamp ON vehicle_locations(vehicle_id, timestamp DESC);