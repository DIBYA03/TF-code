-- +goose Up
/* Linked card table */
CREATE TABLE business_linked_card (
	id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	business_id                 uuid NOT NULL REFERENCES business(id),
	contact_id                  uuid REFERENCES business_contact(id) DEFAULT NULL,
	registered_card_id          text NOT NULL,
	registered_bank_name        text NOT NULL,
	card_number_masked          text NOT NULL,
	card_brand                  text NOT NULL,
	card_type                   text NOT NULL,
	issuer_name                 text NOT NULL,
	fast_funds_enabled          boolean NOT NULL DEFAULT false,
	card_holder_name            text NOT NULL,
	alias                       text DEFAULT NULL,
	account_permission          text DEFAULT 'receiveOnly',
	billing_address             jsonb DEFAULT NULL,
	verified                    boolean NOT NULL DEFAULT false,
	deactivated                 timestamp with time zone DEFAULT NULL,
	created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_linked_card_business_id_fkey ON business_linked_card (business_id);
CREATE INDEX business_linked_card_contact_id_fkey ON business_linked_card (contact_id);
CREATE UNIQUE INDEX business_linked_card_reg_card_id_reg_bank_name_idx
ON business_linked_card (registered_card_id, registered_bank_name);

CREATE TRIGGER update_business_linked_card_modified BEFORE UPDATE ON business_linked_card FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_linked_card_modified on business_linked_card;
DROP INDEX business_linked_card_reg_card_id_reg_bank_name_idx;
DROP INDEX business_linked_card_contact_id_fkey;
DROP INDEX business_linked_card_business_id_fkey;
DROP TABLE business_linked_card;
