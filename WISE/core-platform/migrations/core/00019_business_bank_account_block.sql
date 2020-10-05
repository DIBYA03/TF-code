-- +goose Up
CREATE TABLE business_bank_account_block (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	account_id                  uuid REFERENCES business_bank_account (id),
	reason                      text NOT NULL DEFAULT '',
	block_id                    text NOT NULL,
	block_type                  text NOT NULL,
	deactivated                 timestamp with time zone DEFAULT NULL,
	originated_from             text NOT NULL DEFAULT 'csp',
	created                     timestamp with time zone  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_bank_account_block_account_id_fk ON business_bank_account_block(account_id);
CREATE INDEX business_bank_account_block_account_id_block_id_idx ON business_bank_account_block(account_id,block_id);

-- +goose Down
DROP INDEX business_bank_account_block_account_id_block_id_idx;
DROP INDEX business_bank_account_block_account_id_fk;
DROP TABLE business_bank_account_block;
