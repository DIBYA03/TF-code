-- +goose Up
/* Wise user/business consumer property table maps email/address to ID */
CREATE TABLE consumer_property (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                     uuid NOT NULL REFERENCES consumer (id),
	property_type                   text NOT NULL,
	bank_id                         text NOT NULL,
	bank_name                       text NOT NULL,
	property_value                  jsonb NOT NULL DEFAULT '{}'::jsonb,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

/*
 * Index partner bank consumer id
 */
CREATE INDEX consumer_property_consumer_id_fk ON consumer_property (consumer_id);

/*
 * Index partner bank consumer id and type (unique index)
 */
CREATE UNIQUE INDEX consumer_property_consumer_id_property_type_idx ON consumer_property(consumer_id, property_type);

/*
 * Index by partner property id
 */
CREATE INDEX consumer_property_bank_id_bank_name_idx ON consumer_property (bank_id, bank_name);

CREATE TRIGGER update_consumer_property_modified BEFORE UPDATE ON consumer_property FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_property_modified ON consumer_property;
DROP INDEX consumer_property_bank_id_bank_name_idx;
DROP INDEX consumer_property_consumer_id_property_type_idx;
DROP INDEX consumer_property_consumer_id_fk;
DROP TABLE consumer_property;
