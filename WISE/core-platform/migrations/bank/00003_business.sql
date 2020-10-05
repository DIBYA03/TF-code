-- +goose Up
/*
 * Wise business entity and partner bank business association table
 */
CREATE TABLE business (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                     uuid NOT NULL,
	bank_id                         text NOT NULL,
	bank_name                       text NOT NULL,
	bank_extra                      jsonb NOT NULL DEFAULT '{}'::jsonb,
	kyc_status                      text NOT NULL DEFAULT 'notStarted',
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

/*
 * Business id bank name index
 */
CREATE INDEX business_business_id_bank_name_idx ON business (business_id, bank_name);

/*
 * Bank id bank name key
 */
CREATE INDEX business_bank_id_bank_name_idx ON business (bank_id, bank_name);

CREATE TRIGGER update_business_modified BEFORE UPDATE ON business FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_modified on business;
DROP INDEX business_bank_id_bank_name_idx;
DROP INDEX business_business_id_bank_name_idx;
DROP TABLE business;
