-- A microcosm of the Keycloak database schema for testing

CREATE TABLE IF NOT EXISTS realms (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_entity (
    id VARCHAR(36) PRIMARY KEY,
    realm_id VARCHAR(36),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    enabled BOOLEAN DEFAULT TRUE
);

INSERT INTO realms (id, name) VALUES 
    ('realm-1', 'master')
ON CONFLICT (name) DO NOTHING;

INSERT INTO user_entity (id, realm_id, username, email, first_name, last_name)
VALUES 
    ('user-1', 'realm-1', 'johndoe', 'john@example.com', 'John', 'Doe'),
    ('user-2', 'realm-1', 'janesmith', 'jane@example.com', 'Jane', 'Smith')
ON CONFLICT (email) DO NOTHING;
