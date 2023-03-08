-- Active: 1678179752474@@127.0.0.1@5432@mydb
CREATE TABLE IF NOT EXISTS polls (
	poll_id SERIAL PRIMARY KEY,
	survey_id BIGINT UNIQUE NOT NULL,
	pre_set_values JSON NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	user_id SERIAL PRIMARY KEY,
	email VARCHAR UNIQUE NOT NULL,
	secret VARCHAR NOT NULL,
	created_at TIMESTAMPZ NOT NULL
);