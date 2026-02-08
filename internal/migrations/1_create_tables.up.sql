CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- This migration is intended to be for Postgres. Make another one for your db if the SQL is not accepted by your db.
CREATE EXTENSION IF NOT EXISTS citext;

DO $$
BEGIN
    -- Check if the domain 'email' does not exist
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type 
        WHERE typname = 'email'
    ) THEN
        -- Create the domain if it doesn't exist
        CREATE DOMAIN email AS citext CHECK (
            value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$'
        );
    END IF;
END
$$;



CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4(),
    first_name varchar(255),
    last_name varchar(255),
    email email,
    password varchar(255),
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT email_unique UNIQUE (email)
);