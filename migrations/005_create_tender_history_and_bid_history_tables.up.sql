
CREATE TABLE tender_history (
    id VARCHAR PRIMARY KEY,
    tender_id VARCHAR NOT NULL, 
    name VARCHAR(100),
    description TEXT,
    service_type VARCHAR,
    status VARCHAR,
    organization_id VARCHAR,
    creator_username VARCHAR(50),
    version INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE bid_history (
    id VARCHAR PRIMARY KEY,
    bid_id VARCHAR NOT NULL,
    name VARCHAR(100),
    description TEXT,
    status VARCHAR,
    tender_id VARCHAR,
    author_type VARCHAR,
    author_id VARCHAR,
    creator_username VARCHAR(50),
    version INTEGER,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    valid_from TIMESTAMP,
    valid_to TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

