-- +goose Up
ALTER TABLE business_bank_card ADD COLUMN card_number_alias text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE business_bank_card DROP COLUMN card_number_alias;