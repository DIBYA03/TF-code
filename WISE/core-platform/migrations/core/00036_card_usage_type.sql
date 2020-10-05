-- +goose Up
ALTER TABLE business_linked_card ADD COLUMN usage_type TEXT DEFAULT NULL;
ALTER TABLE business_linked_card ADD COLUMN card_number_hashed TEXT DEFAULT NULL;
UPDATE business_linked_card SET usage_type = 'contact';

-- +goose Down
ALTER TABLE business_linked_card DROP COLUMN card_number_hashed;
ALTER TABLE business_linked_card DROP COLUMN usage_type;

