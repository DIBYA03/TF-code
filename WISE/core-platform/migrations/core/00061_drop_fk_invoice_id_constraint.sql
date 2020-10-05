-- +goose Up
ALTER TABLE business_receipt DROP CONSTRAINT IF EXISTS business_receipt_invoice_id_fk;
ALTER TABLE business_receipt ALTER COLUMN invoice_id DROP NOT NULL;

ALTER TABLE business_receipt DROP CONSTRAINT IF EXISTS  business_receipt_request_id_fk;
ALTER TABLE business_receipt ALTER COLUMN request_id DROP NOT NULL;

ALTER TABLE business_receipt ADD COLUMN invoice_id_v2 UUID;
ALTER TABLE business_receipt ADD CONSTRAINT business_receipt_invoice_v2_unq UNIQUE (invoice_id_v2);
-- +goose Down
ALTER TABLE business_receipt ADD CONSTRAINT business_receipt_invoice_id_fk FOREIGN KEY (invoice_id) REFERENCES business_invoice (id);
ALTER TABLE business_receipt ADD CONSTRAINT business_receipt_request_id_fk FOREIGN KEY (reques_id) REFERENCES business_money_request (id);
ALTER TABLE business_receipt ALTER COLUMN invoice_id ADD NOT NULL;
ALTER TABLE business_receipt ALTER COLUMN request_id ADD NOT NULL;

ALTER TABLE business_receipt DROP COLUMN invoice_id_v2;