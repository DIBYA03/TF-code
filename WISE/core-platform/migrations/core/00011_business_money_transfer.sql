-- +goose Up
/* Business money transfer table */
CREATE TABLE business_money_transfer (
	id                      	uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	created_user_id           	uuid NOT NULL REFERENCES wise_user (id),
	business_id			uuid NOT NULL REFERENCES business (id),
	contact_id			uuid REFERENCES business_contact (id) DEFAULT NULL,
	bank_name               	text NOT NULL,
	bank_transfer_id        	text NOT NULL,
        bank_extra              	jsonb NOT NULL DEFAULT '{}'::jsonb,
        source_account_id         	uuid NOT NULL,
        source_type       	        text NOT NULL,
        dest_account_id           	uuid NOT NULL,
        dest_type         	        text NOT NULL,
        amount                  	double precision NOT NULL,
        currency                	text NOT NULL,
	notes                	        text DEFAULT NULL,
	status                  	text NOT NULL,
        send_email			boolean NOT NULL DEFAULT false,
	posted_transaction_id           uuid DEFAULT NULL,
	created                 	timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        modified                	timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_money_transfer_created_user_id_fkey ON business_money_transfer (created_user_id);
CREATE INDEX business_money_transfer_business_id_fkey ON business_money_transfer (business_id);
CREATE INDEX business_money_transfer_contact_id_fkey ON business_money_transfer (contact_id);
CREATE INDEX business_money_transfer_source_account_id_idx ON business_money_transfer (source_account_id);
CREATE INDEX business_money_transfer_dest_account_id_idx ON business_money_transfer (dest_account_id);
CREATE INDEX business_money_transfer_posted_transaction_id_idx ON business_money_transfer (posted_transaction_id);
CREATE UNIQUE INDEX business_money_transfer_bank_transfer_id_bank_name_idx ON business_money_transfer (bank_transfer_id, bank_name);

CREATE TRIGGER update_business_money_transfer_modified BEFORE UPDATE ON business_money_transfer FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_money_transfer_modified on business_money_transfer;
DROP INDEX business_money_transfer_posted_transaction_id_idx;
DROP INDEX business_money_transfer_created_user_id_fkey;
DROP INDEX business_money_transfer_business_id_fkey;
DROP INDEX business_money_transfer_contact_id_fkey;
DROP INDEX business_money_transfer_source_account_id_idx;
DROP INDEX business_money_transfer_dest_account_id_idx;
DROP INDEX business_money_transfer_bank_transfer_id_bank_name_idx;
DROP TABLE business_money_transfer;
