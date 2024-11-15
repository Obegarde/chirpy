-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	email TEXT NOT NULL UNIQUE,
	is_chirpy_red BOOLEAN DEFAULT FALSE,
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

CREATE TABLE refresh_tokens(
	token text PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
DROP TABLE chirps;
DROP TABLE users;
