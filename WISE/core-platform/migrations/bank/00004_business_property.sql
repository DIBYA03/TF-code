-- +goose Up
/* Wise business property table maps email/address to ID */
CREATE TABLE business_property (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                     uuid NOT NULL REFERENCES business (id),
	property_type                   text NOT NULL,
	bank_id                         text NOT NULL,
	bank_name                       text NOT NULL,
	property_value                  jsonb NOT NULL DEFAULT '{}'::jsonb,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

/*
 * Index partner bank business id
 */
CREATE INDEX business_property_business_id_fk ON business_property (business_id);

/*
 * Index partner bank business id and type (unique index)
 */
CREATE UNIQUE INDEX business_property_business_id_property_type_idx ON business_property(business_id, property_type);

/*
 * Index by partner property id
 */
CREATE INDEX business_property_bank_id_bank_name_idx ON business_property (bank_id, bank_name);

CREATE TRIGGER update_business_property_modified BEFORE UPDATE ON business_property FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_property_modified on business_property;
DROP INDEX business_property_bank_id_bank_name_idx;
DROP INDEX business_property_business_id_property_type_idx;
DROP INDEX business_property_business_id_fk;
DROP TABLE business_property;
