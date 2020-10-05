-- +goose Up

/* Transaction table */
ALTER TABLE business_transaction ADD COLUMN amount_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_transaction SET amount_fd = CAST(amount AS DECIMAL(19,4));

/* Card transaction table */
ALTER TABLE business_card_transaction ADD COLUMN auth_amount_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_card_transaction SET auth_amount_fd = CAST(auth_amount AS DECIMAL(19,4)); 

ALTER TABLE business_card_transaction ADD COLUMN local_amount_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_card_transaction SET local_amount_fd = CAST(local_amount AS DECIMAL(19,4)); 

/* Hold transaction */
ALTER TABLE business_hold_transaction ADD COLUMN amount_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_hold_transaction SET amount_fd = CAST(amount AS DECIMAL(19,4));

/* Account daily balance */
ALTER TABLE business_account_daily_balance ADD COLUMN posted_balance_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_account_daily_balance SET posted_balance_fd = CAST(posted_balance AS DECIMAL(19,4));

ALTER TABLE business_account_daily_balance ADD COLUMN apy_bps INT NOT NULL DEFAULT 0;
UPDATE business_account_daily_balance SET apy_bps = CAST(apy*10000 AS INT);

ALTER TABLE business_account_daily_balance RENAME COLUMN money_credited TO amount_credited;
ALTER TABLE business_account_daily_balance ADD COLUMN amount_credited_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_account_daily_balance SET amount_credited_fd = CAST(amount_credited AS DECIMAL(19,4));

ALTER TABLE business_account_daily_balance RENAME COLUMN money_debited TO amount_debited;
ALTER TABLE business_account_daily_balance ADD COLUMN amount_debited_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_account_daily_balance SET amount_debited_fd = CAST(amount_debited AS DECIMAL(19,4));

/* Daily transaction stats */
ALTER TABLE business_daily_transaction_stats RENAME COLUMN money_requested TO amount_requested;
ALTER TABLE business_daily_transaction_stats ADD COLUMN amount_requested_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_daily_transaction_stats SET amount_requested_fd = CAST(amount_requested AS DECIMAL(19,4));

ALTER TABLE business_daily_transaction_stats RENAME COLUMN money_paid TO amount_paid;
ALTER TABLE business_daily_transaction_stats ADD COLUMN amount_paid_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_daily_transaction_stats SET amount_paid_fd = CAST(amount_paid AS DECIMAL(19,4));

ALTER TABLE business_daily_transaction_stats RENAME COLUMN money_sent TO amount_sent;
ALTER TABLE business_daily_transaction_stats ADD COLUMN amount_sent_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_daily_transaction_stats SET amount_sent_fd = CAST(amount_sent AS DECIMAL(19,4));

ALTER TABLE business_daily_transaction_stats RENAME COLUMN money_credited TO amount_credited;
ALTER TABLE business_daily_transaction_stats ADD COLUMN amount_credited_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_daily_transaction_stats SET amount_credited_fd = CAST(amount_credited AS DECIMAL(19,4));

ALTER TABLE business_daily_transaction_stats RENAME COLUMN money_debited TO amount_debited;
ALTER TABLE business_daily_transaction_stats ADD COLUMN amount_debited_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_daily_transaction_stats SET amount_debited_fd = CAST(amount_debited AS DECIMAL(19,4));

/* Alloy results */
ALTER TABLE business_transaction_alloy_result ADD COLUMN amount_fd DECIMAL(19,4) NOT NULL DEFAULT 0;
UPDATE business_transaction_alloy_result SET amount_fd = CAST(amount AS DECIMAL(19,4));

-- +goose Down
ALTER TABLE business_transaction_alloy_result DROP COLUMN amount_fd;

ALTER TABLE business_daily_transaction_stats DROP COLUMN amount_debited_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_debited TO money_debited;

ALTER TABLE business_daily_transaction_stats DROP COLUMN amount_credited_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_credited TO money_credited;

ALTER TABLE business_daily_transaction_stats DROP COLUMN amount_sent_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_sent TO money_sent;

ALTER TABLE business_daily_transaction_stats DROP COLUMN amount_paid_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_paid TO money_paid;

ALTER TABLE business_daily_transaction_stats DROP COLUMN amount_requested_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_requested TO money_requested;

ALTER TABLE business_account_daily_balance DROP COLUMN amount_debited_fd;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_debited TO money_debited;

ALTER TABLE business_account_daily_balance DROP COLUMN amount_credited_fd;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_credited TO money_credited;

ALTER TABLE business_account_daily_balance DROP COLUMN apy_bps;
ALTER TABLE business_account_daily_balance DROP COLUMN posted_balance_fd;

ALTER TABLE business_hold_transaction DROP COLUMN amount_fd;

ALTER TABLE business_card_transaction DROP COLUMN local_amount_fd;
ALTER TABLE business_card_transaction DROP COLUMN auth_amount_fd;

ALTER TABLE business_transaction DROP COLUMN amount_fd;
