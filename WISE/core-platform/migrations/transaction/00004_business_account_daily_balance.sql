-- +goose Up
/* Business bank account table */
CREATE TABLE business_account_daily_balance (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	account_id                      uuid NOT NULL,
	business_id                     uuid NOT NULL,
	posted_balance                  double precision NOT NULL,
	currency                        text NOT NULL,
	apy                             double precision NOT NULL DEFAULT 0,
	recorded_date                   date NOT NULL DEFAULT CURRENT_DATE,
	created                         timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_account_daily_balance_business_id_account_id_recorded_date_idx ON business_account_daily_balance (business_id, account_id, recorded_date);

-- +goose Down
DROP INDEX business_account_daily_balance_business_id_account_id_recorded_date_idx;
DROP TABLE business_account_daily_balance;
