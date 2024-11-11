-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	email TEXT NOT NULL UNIQUE,
	hashed_password TEXT NOT NULL DEFAULT 'unset'

);
CREATE TABLE chirps(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	body TEXT NOT NULL,
	user_id UUID NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
DROP TABLE users;
