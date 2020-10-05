-- +goose Up
/* Business bank account table */
CREATE TABLE business_bank_account (
	id                      	uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
        account_holder_id         	uuid NOT NULL REFERENCES wise_user (id),
	business_id			uuid NOT NULL REFERENCES business (id),
	bank_name               	text NOT NULL,
	bank_account_id         	text NOT NULL,
        bank_extra              	jsonb NOT NULL DEFAULT '{}'::jsonb,
        account_type             	text NOT NULL,
        account_status           	text NOT NULL,
        account_number           	text NOT NULL,
        routing_number           	text NOT NULL,
	wire_routing                    text DEFAULT NULL,
        alias                   	text DEFAULT NULL,
        available_balance        	double precision NOT NULL,
        posted_balance           	double precision NOT NULL,
        currency                	text NOT NULL,
	interest_ytd                    double precision DEFAULT 0,
	monthly_cycle_start_day         int DEFAULT 1,
        opened                  	timestamp with time zone NOT NULL,
        created                 	timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        modified                	timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_bank_account_account_holder_id_fkey ON business_bank_account (account_holder_id);
CREATE INDEX business_bank_account_business_id_fkey ON business_bank_account (business_id);
CREATE UNIQUE INDEX business_bank_account_bank_account_id_bank_name_idx ON business_bank_account (bank_account_id, bank_name);

CREATE TRIGGER update_business_bank_account_modified BEFORE UPDATE ON business_bank_account FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_bank_account_modified on business_bank_account;
DROP INDEX business_bank_account_business_id_fkey;
DROP INDEX business_bank_account_account_holder_id_fkey;
DROP TABLE business_bank_account;
