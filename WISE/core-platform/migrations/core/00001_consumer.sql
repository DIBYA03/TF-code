-- +goose Up
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

/* Consumer table */
CREATE TABLE consumer (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	first_name                      text NOT NULL DEFAULT '',
	middle_name                     text NOT NULL DEFAULT '',
	last_name                       text NOT NULL DEFAULT '',
	email                           text DEFAULT NULL,
	phone                           text DEFAULT NULL,
	date_of_birth                   date DEFAULT NULL,
	tax_id                          text DEFAULT NULL,
	tax_id_type                     text DEFAULT NULL,
	kyc_status                      text NOT NULL DEFAULT 'notStarted',
	legal_address                   jsonb DEFAULT NULL,
	mailing_address                 jsonb DEFAULT NULL,
	work_address                    jsonb DEFAULT NULL,
	residency                       jsonb DEFAULT NULL,
	citizenship_countries           jsonb NOT NULL DEFAULT '[]'::jsonb,
	occupation                      text DEFAULT NULL,
	income_type                     jsonb NOT NULL DEFAULT '[]'::jsonb,
	activity_type                   jsonb NOT NULL DEFAULT
	                                '['
	                                        '"check",'
	                                        '"domesticWireTransfer",'
	                                        '"internationalWireTransfer",'
	                                        '"domesticACH",'
	                                        '"internationalACH"'
	                                ']'::jsonb,
	is_restricted                   boolean NOT NULL DEFAULT false,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_consumer_modified BEFORE UPDATE ON consumer FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_modified on consumer;
DROP TABLE consumer;

DROP FUNCTION IF EXISTS update_modified_column;
