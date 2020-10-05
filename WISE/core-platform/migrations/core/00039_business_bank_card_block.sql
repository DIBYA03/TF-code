-- +goose Up
CREATE TABLE business_bank_card_block (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	card_id                     uuid REFERENCES business_bank_card (id),
	block_id                    text NOT NULL,
	reason                      text NOT NULL DEFAULT '',
	originated_from             text NOT NULL,
	block_status                text NOT NULL,
	created                     timestamp with time zone  NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified 	            timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_bank_card_block_card_id_fk ON business_bank_card_block(card_id);
CREATE INDEX business_bank_card_block_card_id_block_id_idx ON business_bank_card_block(card_id,block_id);

CREATE TRIGGER update_business_bank_card_block_modified BEFORE UPDATE ON business_bank_card_block FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_bank_card_block_modified on business_bank_card_block;
DROP INDEX business_bank_card_block_card_id_block_id_idx;
DROP INDEX business_bank_card_block_card_id_fk;
DROP TABLE business_bank_card_block;
