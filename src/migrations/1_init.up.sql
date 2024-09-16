CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    user_id INT REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type VARCHAR(50),
    status VARCHAR(50) CHECK (status IN ('CREATED', 'PUBLISHED', 'CLOSED')),
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    version INT
);

CREATE TABLE IF NOT EXISTS bid (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(50) CHECK (status IN ('CREATED', 'PUBLISHED', 'CANCELED', 'SUBMITTED', 'DECLINED')),
    tender_id INT REFERENCES tender(id) ON DELETE CASCADE,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    creator_username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    version INT
);

CREATE TABLE IF NOT EXISTS tender_history (
    id SERIAL PRIMARY KEY,
    tender_id INT NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    service_type VARCHAR(50),
    status VARCHAR(50),
    organization_id INT,
    creator_username VARCHAR(50),
    version INT
);

CREATE TABLE IF NOT EXISTS bid_history (
    id SERIAL PRIMARY KEY,
    bid_id INT NOT NULL REFERENCES bid(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(50),
    tender_id INT,
    organization_id INT,
    creator_username VARCHAR(50),
    version INT
);

CREATE TABLE IF NOT EXISTS bid_decisions (
    id SERIAL PRIMARY KEY,
    bid_id INT NOT NULL REFERENCES bid(id) ON DELETE CASCADE,
    username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    decision_status VARCHAR(50) CHECK (decision_status IN ('SUBMITTED', 'DECLINED')),
    decision_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    bid_id INT NOT NULL REFERENCES bid(id) ON DELETE CASCADE,
    username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE,
    organization_id INT REFERENCES organization(id) ON DELETE CASCADE,
    comment TEXT NOT NULL
);

INSERT INTO employee (username, first_name, last_name) VALUES
('user1', 'John', 'Doe'),
('user2', 'Jane', 'Smith'),
('user3', 'Kate', 'Jones'),
('user4', 'Mary', 'Smith'),
('user5', 'Lore', 'Simpson'),
('user6', 'Sandra', 'Skale'),
('simpleUser1', 'Eliza', 'Ted'),
('simpleUser2', 'Sofi', 'Sun');

INSERT INTO organization (name, description, type) VALUES
('Organization 1', 'Description for Organization 1', 'LLC'),
('Organization 2', 'Description for Organization 2', 'JSC');

INSERT INTO organization_responsible (organization_id, user_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(2, 4),
(2, 5),
(2, 6);