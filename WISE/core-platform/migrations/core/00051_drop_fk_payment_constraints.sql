-- +goose Up
ALTER TABLE business_money_request_payment DROP CONSTRAINT business_money_request_payment_linked_bank_account_id_fkey;

-- +goose Down
ALTER TABLE business_money_request_payment ADD CONSTRAINT business_money_request_payment_linked_bank_account_id_fkey FOREIGN KEY (linked_bank_account_id) REFERENCES business_linked_bank_account (id);
