-- +goose Up
DROP INDEX IF EXISTS business_daily_transaction_stats_business_id_recorded_date_idx;
CREATE UNIQUE INDEX business_daily_transaction_stats_business_id_recorded_date_idx ON business_daily_transaction_stats (business_id, recorded_date);

-- +goose Down
DROP INDEX IF EXISTS business_daily_transaction_stats_business_id_recorded_date_idx;
