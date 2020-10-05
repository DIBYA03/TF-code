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

/* Wise user table */
CREATE TABLE identity (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	provider_id                     text NOT NULL,
	provider_name                   text NOT NULL,
	provider_source                 text NOT NULL,
	phone                           text UNIQUE NOT NULL,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX identity_provider_id_provider_name_provider_source_idx ON identity (provider_id, provider_name, provider_source);
CREATE TRIGGER update_identity_modified BEFORE UPDATE ON identity FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_identity_modified on identity;
DROP INDEX identity_provider_id_provider_name_provider_source_idx;
DROP TABLE identity;

DROP FUNCTION IF EXISTS update_modified_column;
