-- +goose Up
/* Linked account table */
CREATE TABLE business_linked_payee (
    id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    business_id                 uuid NOT NULL REFERENCES business(id),
    contact_id                  uuid REFERENCES business_contact(id) DEFAULT NULL,
    address_id                  uuid DEFAULT NULL,
    bank_payee_id               text NOT NULL,
    bank_name                   text NOT NULL,
    account_holder_name         text NOT NULL,
    payee_name                  text NOT NULL,
    status                      text NOT NULL,
    created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_linked_payee_business_id_idx ON business_linked_payee (business_id);
CREATE INDEX business_linked_payee_contact_id_idx ON business_linked_payee (contact_id);
CREATE INDEX business_linked_payee_address_id_idx ON business_linked_payee (address_id);
CREATE INDEX business_linked_payee_bank_payee_id_bank_name ON business_linked_payee (bank_payee_id, bank_name);

CREATE TRIGGER update_business_linked_payee_modified BEFORE UPDATE ON business_linked_payee FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_linked_payee_modified on business_linked_payee;
DROP INDEX business_linked_payee_bank_payee_id_bank_name;
DROP INDEX business_linked_payee_business_id_idx;
DROP INDEX business_linked_payee_contact_id_idx;
DROP INDEX business_linked_payee_address_id_idx;
DROP TABLE business_linked_payee;
