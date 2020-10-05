-- +goose Up
/* Add UUID extension  */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
/* Add modified update function */
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- +goose StatementEnd

CREATE TABLE csp_user(
	id           uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	first_name   text NOT NULL DEFAULT '',
	middle_name  text NOT NULL DEFAULT '',
	last_name    text NOT NULL DEFAULT '',
	email         text DEFAULT NULL,
	email_verified boolean NOT NULL DEFAULT false,
	phone          text UNIQUE NOT NULL,
	phone_verified boolean NOT NULL DEFAULT false,
	created        timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified       timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE csp_user;
