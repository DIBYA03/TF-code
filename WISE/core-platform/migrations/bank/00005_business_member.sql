-- +goose Up
/* Wise user/business entity property table maps email/address to ID */
CREATE TABLE business_member (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                     uuid NOT NULL REFERENCES consumer (id),
	business_id                     uuid NOT NULL REFERENCES business (id),
	bank_id                         text DEFAULT NULL,
	bank_name                       text NOT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

/*
 * Consumer id index
 */
CREATE INDEX business_member_consumer_id_fk ON business_member (consumer_id);

/*
 * Business id index
 */
CREATE INDEX business_member_business_id_fk ON business_member (business_id);

/*
 * Index bank id bank name
 */
CREATE UNIQUE INDEX business_member_bank_id_bank_name_idx ON business_member(bank_id, bank_name);

/* 
 *
 */
CREATE UNIQUE INDEX business_member_business_id_consumer_id_idx ON business_member(business_id, consumer_id);

CREATE TRIGGER update_business_member_modified BEFORE UPDATE ON business_member FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_member_modified on business_member;
DROP INDEX business_member_business_id_consumer_id_idx;
DROP INDEX business_member_bank_id_bank_name_idx;
DROP INDEX business_member_business_id_fk;
DROP INDEX business_member_consumer_id_fk;
DROP TABLE business_member;
