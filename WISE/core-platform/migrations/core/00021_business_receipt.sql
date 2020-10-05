-- +goose Up
CREATE TABLE business_receipt ( 
    id                         uuid PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    request_id                 uuid REFERENCES business_money_request (id), 
    invoice_id                 uuid REFERENCES business_invoice (id), 
    created_user_id            uuid NOT NULL REFERENCES wise_user (id), 
    business_id                uuid NOT NULL REFERENCES business (id), 
    contact_id                 uuid REFERENCES business_contact (id) DEFAULT NULL, 
    receipt_number             TEXT NOT NULL, 
    storage_key                TEXT NOT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_receipt_request_id_fk ON business_receipt(request_id);
CREATE INDEX business_receipt_created_user_id_fk ON business_receipt(created_user_id);
CREATE INDEX business_receipt_business_id_fk ON business_receipt(business_id);
CREATE INDEX business_receipt_contact_id_fk ON business_receipt(contact_id);
CREATE INDEX business_receipt_invoice_id_fk ON business_receipt(invoice_id);

-- +goose Down
DROP INDEX business_receipt_request_id_fk;
DROP INDEX business_receipt_created_user_id_fk;
DROP INDEX business_receipt_business_id_fk;
DROP INDEX business_receipt_contact_id_fk;
DROP INDEX business_receipt_invoice_id_fk;
DROP TABLE business_receipt;
