-- +goose Up
ALTER TABLE business_money_request_payment DROP CONSTRAINT business_money_request_payment_request_id_fkey;
ALTER TABLE business_money_request_payment ALTER COLUMN request_id DROP NOT NULL;
-- +goose Down
ALTER TABLE business_money_request_payment ADD CONSTRAINT business_money_request_payment_request_id_fkey FOREIGN KEY (request_id) REFERENCES business_money_request (id);
ALTER TABLE business_money_request_payment ALTER COLUMN request_id ADD NOT NULL;