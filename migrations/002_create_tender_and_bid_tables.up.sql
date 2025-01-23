
CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
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
    'CREATED',
    'PUBLISHED',
    'CANCELED',
    'APPROVED',
    'REJECTED'
);

-- Создание ENUM типа для типов авторов предложений
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