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

CREATE TABLE IF NOT EXISTS books (
    /* Use uuid if desired more security, but be aware of performance reduction with big DBs*/
    id serial primary key,
    title varchar(255),
    author varchar(255),
    cover_url varchar(255),
    post_url varchar(255),
    created_at timestamp,
    updated_at timestamp,    
    CONSTRAINT id_unique UNIQUE (id)
);


