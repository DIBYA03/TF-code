-- +goose Up
/* PLATFORM-319: currency is already being tracked and user id needs to be added for analytics purpose */
ALTER TABLE business_card_transaction DROP COLUMN currency;
ALTER TABLE business_card_transaction ADD COLUMN cardholder_id text DEFAULT NULL;

-- +goose Down
ALTER TABLE business_card_transaction DROP COLUMN cardholder_id;
ALTER TABLE business_card_transaction ADD COLUMN currency text DEFAULT NULL;
