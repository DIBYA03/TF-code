-- +goose Up
/* Business bank account table */
CREATE TABLE business_account_monthly_interest (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	account_id                      uuid NOT NULL,
	business_id                     uuid NOT NULL,
	avg_posted_balance              DECIMAL(19,4) NOT NULL,
	days                            int NOT NULL,
	interest_amount                 DECIMAL(19,4) NOT NULL,
	interest_payout                 DECIMAL(19,2) NOT NULL,
	currency                        text NOT NULL,
	apr                             int NOT NULL DEFAULT 0,
	start_date                      date NOT NULL,
	end_date                        date NOT NULL,
	recorded_date                   date NOT NULL DEFAULT CURRENT_DATE,
	created                         timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX business_account_monthly_interest_business_id_account_id_recorded_date_idx ON business_account_monthly_interest (business_id, account_id, recorded_date);

ALTER TABLE business_account_daily_balance DROP COLUMN apy;
ALTER TABLE business_account_daily_balance RENAME COLUMN apy_bps TO apr;

-- +goose Down
ALTER TABLE business_account_daily_balance RENAME COLUMN apr TO apy_bps;
ALTER TABLE business_account_daily_balance ADD COLUMN apy double precision NOT NULL DEFAULT 0;

DROP INDEX business_account_monthly_interest_business_id_account_id_recorded_date_idx;
DROP TABLE business_account_monthly_interest;
