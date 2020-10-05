-- +goose Up
ALTER TABLE business_bank_card_block ADD COLUMN block_date timestamp with time zone DEFAULT NULL;

-- +goose Down
ALTER TABLE business_bank_card_block DROP COLUMN block_date;