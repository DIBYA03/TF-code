-- +goose Up
ALTER TABLE money_transfer_request DROP CONSTRAINT money_transfer_request_money_transfer_id_fkey;

-- +goose Down
ALTER TABLE money_transfer_request ADD CONSTRAINT money_transfer_request_money_transfer_id_fkey FOREIGN KEY (money_transfer_id) REFERENCES business_money_transfer (id);
