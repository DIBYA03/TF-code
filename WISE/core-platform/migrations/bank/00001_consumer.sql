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

/*
 * Partner bank consumer table
 */
CREATE TABLE consumer (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                     uuid NOT NULL,
	bank_id                         text NOT NULL,
	bank_name                       text NOT NULL,
	bank_extra                      jsonb NOT NULL DEFAULT '{}'::jsonb,
	kyc_status                      text NOT NULL DEFAULT 'notStarted',
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

/*
 * Consumer id unique index
 */
CREATE UNIQUE INDEX consumer_consumer_id_bank_name_idx ON consumer (consumer_id, bank_name);

/*
 * Bank id unique index
 */
CREATE UNIQUE INDEX consumer_bank_id_bank_name_idx ON consumer (bank_id, bank_name);

CREATE TRIGGER update_consumer_modified BEFORE UPDATE ON consumer FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_modified on consumer;
DROP INDEX consumer_bank_id_bank_name_idx;
DROP INDEX consumer_consumer_id_bank_name_idx;
DROP TABLE consumer;

DROP FUNCTION IF EXISTS update_modified_column;
