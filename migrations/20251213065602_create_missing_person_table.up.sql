CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE missing_persons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    age INT,
    description TEXT NOT NULL,
    last_seen VARCHAR(255) NOT NULL,
    contact VARCHAR(100) NOT NULL,
    photo_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
