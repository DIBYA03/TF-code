-- +goose Up
/* Linked account table */
CREATE TABLE business_linked_bank_account (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	user_id                     uuid REFERENCES wise_user (id) DEFAULT NULL,
	business_id                 uuid NOT NULL REFERENCES business(id),
	business_bank_account_id    uuid REFERENCES business_bank_account(id) DEFAULT NULL,
	contact_id                  uuid REFERENCES business_contact(id) DEFAULT NULL,
	registered_account_id       text NOT NULL,
	registered_bank_name        text NOT NULL,
	account_holder_name         text NOT NULL,
	account_name                text DEFAULT NULL,
	currency                    text DEFAULT NULL,
	account_type                text NOT NULL,
	account_number              text NOT NULL,
	bank_name                   text DEFAULT NULL,
	routing_number              text NOT NULL,
	wire_routing                text DEFAULT NULL,
	source_account_id           text DEFAULT NULL,
	source_id                   text DEFAULT NULL,
	source_name                 text DEFAULT NULL,
	account_permission          text DEFAULT 'receiveOnly',
	alias                       text DEFAULT NULL,
	verified                    boolean NOT NULL DEFAULT false,
	deactivated	            timestamp with time zone DEFAULT NULL,
	created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_linked_bank_account_user_id_fkey ON business_linked_bank_account (user_id);
CREATE INDEX business_linked_bank_account_business_id_fkey ON business_linked_bank_account (business_id);
CREATE INDEX business_linked_bank_account_business_bank_account_id_fkey ON business_linked_bank_account (business_bank_account_id);
CREATE INDEX business_linked_bank_account_contact_id_fkey ON business_linked_bank_account (contact_id);
CREATE INDEX business_linked_bank_acc_reg_bank_name_acc_num_routing_num_idx
ON business_linked_bank_account (business_id, account_number, routing_number);
CREATE UNIQUE INDEX business_linked_bank_account_reg_account_id_reg_bank_name_idx
ON business_linked_bank_account (registered_account_id, registered_bank_name);
CREATE INDEX business_linked_bank_account_source_account_id_source_name_idx
ON business_linked_bank_account (source_account_id, source_name);

CREATE TRIGGER update_business_linked_bank_account_modified BEFORE UPDATE ON business_linked_bank_account FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_linked_bank_account_modified on business_linked_bank_account;
DROP INDEX business_linked_bank_account_user_id_fkey;
DROP INDEX business_linked_bank_account_business_id_fkey;
DROP INDEX business_linked_bank_account_business_bank_account_id_fkey;
DROP INDEX business_linked_bank_account_contact_id_fkey;
DROP INDEX business_linked_bank_acc_reg_bank_name_acc_num_routing_num_idx;
DROP INDEX business_linked_bank_account_reg_account_id_reg_bank_name_idx;
DROP INDEX business_linked_bank_account_source_account_id_source_name_idx;
DROP TABLE business_linked_bank_account;
