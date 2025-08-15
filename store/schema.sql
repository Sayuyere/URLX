-- Postgres table for URL shortener
CREATE TABLE IF NOT EXISTS urls (
    short VARCHAR(32) PRIMARY KEY,
    long TEXT NOT NULL
);
