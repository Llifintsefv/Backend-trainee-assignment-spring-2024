
CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);


CREATE TABLE employee (
    id VARCHAR PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE organization (
    id VARCHAR PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE organization_responsible (
    id VARCHAR PRIMARY KEY,
    organization_id VARCHAR REFERENCES organization(id) ON DELETE CASCADE,
    user_id VARCHAR REFERENCES employee(id) ON DELETE CASCADE,
    UNIQUE (organization_id, user_id)
);