-- +goose Up
ALTER TABLE business_account_daily_balance ADD COLUMN money_credited DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE business_account_daily_balance ADD COLUMN money_debited DOUBLE PRECISION NOT NULL DEFAULT 0;

DROP INDEX IF EXISTS business_account_daily_balance_business_id_account_id_idx;
CREATE UNIQUE INDEX business_account_daily_balance_business_id_account_id_recorded_date_idx ON business_account_daily_balance (business_id, account_id, recorded_date);

-- +goose Down
DROP INDEX IF EXISTS business_account_daily_balance_business_id_account_id_recorded_date_idx;
ALTER TABLE business_account_daily_balance DROP COLUMN money_debited;
ALTER TABLE business_account_daily_balance DROP COLUMN money_credited;
