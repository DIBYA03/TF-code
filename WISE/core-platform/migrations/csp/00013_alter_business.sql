-- +goose Up 
ALTER TABLE business ADD COLUMN promo_funded timestamp with time zone DEFAULT NULL;
ALTER TABLE business ADD COLUMN promo_money_transfer_id text DEFAULT NULL;
ALTER TABLE business ADD COLUMN employee_count bigint DEFAULT 0;
ALTER TABLE business ADD COLUMN location_count bigint DEFAULT 0;
ALTER TABLE business ADD COLUMN customer_type text DEFAULT '';
ALTER TABLE business ADD COLUMN accept_card_payment text DEFAULT '';

-- +goose Down
ALTER TABLE business DROP COLUMN promo_funded;
ALTER TABLE business DROP COLUMN promo_money_transfer_id;
ALTER TABLE business DROP COLUMN employee_count;
ALTER TABLE business DROP COLUMN location_count;
ALTER TABLE business DROP COLUMN customer_type;
ALTER TABLE business DROP COLUMN accept_card_payment;