-- +goose Up
ALTER TABLE business_member ADD COLUMN bank_control_id TEXT DEFAULT NULL;
DROP INDEX business_member_bank_id_bank_name_idx;
CREATE UNIQUE INDEX business_member_bank_control_id_bank_id_bank_name_idx ON business_member(bank_control_id, bank_id, bank_name);

-- +goose Down
DROP INDEX business_member_bank_control_id_bank_id_bank_name_idx;
CREATE UNIQUE INDEX business_member_bank_id_bank_name_idx ON business_member(bank_id, bank_name);
ALTER TABLE business_member DROP COLUMN bank_control_id;
