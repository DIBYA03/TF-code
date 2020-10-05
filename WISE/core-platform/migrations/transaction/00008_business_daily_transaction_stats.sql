-- +goose Up
/* Business daily transaction stats table */
CREATE TABLE business_daily_transaction_stats ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID NOT NULL, 
    money_requested            DOUBLE PRECISION NOT NULL DEFAULT 0, 
    money_paid                 DOUBLE PRECISION NOT NULL DEFAULT 0, 
    money_sent                 DOUBLE PRECISION NOT NULL DEFAULT 0, 
    money_credited             DOUBLE PRECISION NOT NULL DEFAULT 0, 
    money_debited              DOUBLE PRECISION NOT NULL DEFAULT 0, 
    currency                   TEXT NOT NULL, 
    recorded_date              DATE NOT NULL DEFAULT CURRENT_DATE, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_daily_transaction_stats_business_id_idx ON business_daily_transaction_stats(business_id);
CREATE UNIQUE INDEX business_daily_transaction_stats_business_id_recorded_date_idx ON business_account_daily_balance (business_id, recorded_date);

-- +goose Down
DROP INDEX business_daily_transaction_stats_business_id_recorded_date_idx;
DROP INDEX business_daily_transaction_stats_business_id_idx;
DROP TABLE business_daily_transaction_stats;
