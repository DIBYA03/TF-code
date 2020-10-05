-- +goose Up
ALTER TABLE business_money_request_payment ADD COLUMN card_brand TEXT DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN card_number TEXT DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN payment_date timestamp with time zone DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN receipt_id UUID DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN receipt_mode TEXT DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN customer_contact TEXT DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN receipt_token TEXT DEFAULT NULL;
ALTER TABLE business_money_request_payment ADD COLUMN purchase_address jsonb DEFAULT NULL;
ALTER TABLE business_money_request_payment RENAME COLUMN token TO payment_token;

-- +goose Down
ALTER TABLE business_money_request_payment RENAME COLUMN payment_token TO token;
ALTER TABLE business_money_request_payment DROP COLUMN purchase_address;
ALTER TABLE business_money_request_payment DROP COLUMN receipt_token;
ALTER TABLE business_money_request_payment DROP COLUMN customer_contact;
ALTER TABLE business_money_request_payment DROP COLUMN receipt_mode;
ALTER TABLE business_money_request_payment DROP COLUMN receipt_id;
ALTER TABLE business_money_request_payment DROP COLUMN payment_date;
ALTER TABLE business_money_request_payment DROP COLUMN card_number;
ALTER TABLE business_money_request_payment DROP COLUMN card_brand;
