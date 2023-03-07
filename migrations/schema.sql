-- Active: 1678179752474@@127.0.0.1@5432@polls-db@public
CREATE TABLE IF NOT EXISTS polls (
	poll_id TEXT PRIMARY KEY,
	survey_id INT NOT NULL,
	pre_set_values JSON NOT NULL
);