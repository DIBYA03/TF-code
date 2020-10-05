-- +goose Up
CREATE TABLE business_invoice ( 
    id                         uuid PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    request_id                 uuid REFERENCES business_money_request (id), 
    created_user_id            uuid NOT NULL REFERENCES wise_user (id), 
    business_id                uuid NOT NULL REFERENCES business (id), 
    contact_id                 uuid REFERENCES business_contact (id) DEFAULT NULL, 
    invoice_number             TEXT NOT NULL, 
    storage_key                TEXT NOT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

CREATE INDEX business_invoice_request_id_fk ON business_invoice(request_id);
CREATE INDEX business_invoice_created_user_id_fk ON business_invoice(created_user_id);
CREATE INDEX business_invoice_business_id_fk ON business_invoice(business_id);
CREATE INDEX business_invoice_contact_id_fk ON business_invoice(contact_id);

-- +goose Down
DROP INDEX business_invoice_request_id_fk;
DROP INDEX business_invoice_created_user_id_fk;
DROP INDEX business_invoice_business_id_fk;
DROP INDEX business_invoice_contact_id_fk;
DROP TABLE business_invoice;
