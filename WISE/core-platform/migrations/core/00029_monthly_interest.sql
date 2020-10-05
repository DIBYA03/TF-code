-- +goose Up
ALTER TABLE business_money_transfer ADD COLUMN account_monthly_interest_id uuid DEFAULT NULL;
CREATE INDEX business_money_transfer_account_monthly_interest_id_idx ON business_money_transfer (account_monthly_interest_id);

-- +goose Down
DROP INDEX business_money_transfer_account_monthly_interest_id_idx;
ALTER TABLE business_money_transfer DROP COLUMN business_account_monthly_interest_id;
