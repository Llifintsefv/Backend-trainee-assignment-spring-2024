
CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);


CREATE TYPE tender_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);


CREATE TABLE tender (
    id VARCHAR PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type tender_service_type NOT NULL,
    status tender_status NOT NULL,
    organization_id VARCHAR REFERENCES organization(id) NOT NULL,
    creator_username VARCHAR REFERENCES employee(username) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled',
    'Approved',
    'Rejected'
);


CREATE TYPE bid_author_type AS ENUM (
    'Organization',
    'User'
);


CREATE TABLE bid (
    id VARCHAR PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status bid_status NOT NULL,
    tender_id VARCHAR REFERENCES tender(id) NOT NULL,
    author_type bid_author_type NOT NULL,
    author_id VARCHAR NOT NULL, 
    creator_username VARCHAR REFERENCES employee(username) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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
    updated_at TIMESTAMP,
);

CREATE INDEX idx_tender_history_tender_id ON tender_history (tender_id);


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

CREATE INDEX idx_bid_history_bid_id ON bid_history (bid_id);