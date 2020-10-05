-- +goose Up

/* Transaction table */
ALTER TABLE business_transaction RENAME COLUMN amount TO amount_dep;
ALTER TABLE business_transaction ALTER COLUMN amount_dep SET DEFAULT 0;
ALTER TABLE business_transaction RENAME COLUMN amount_fd TO amount;
ALTER TABLE business_transaction ALTER COLUMN amount DROP DEFAULT;

/* Card transaction table */
ALTER TABLE business_card_transaction RENAME COLUMN auth_amount TO auth_amount_dep;
ALTER TABLE business_card_transaction ALTER COLUMN auth_amount_dep SET DEFAULT 0;
ALTER TABLE business_card_transaction RENAME COLUMN auth_amount_fd TO auth_amount;
ALTER TABLE business_card_transaction ALTER COLUMN auth_amount DROP DEFAULT;

ALTER TABLE business_card_transaction RENAME COLUMN local_amount TO local_amount_dep;
ALTER TABLE business_card_transaction ALTER COLUMN local_amount_dep SET DEFAULT 0;
ALTER TABLE business_card_transaction RENAME COLUMN local_amount_fd TO local_amount;
ALTER TABLE business_card_transaction ALTER COLUMN local_amount DROP DEFAULT;

/* Hold transaction */
ALTER TABLE business_hold_transaction RENAME COLUMN amount TO amount_dep;
ALTER TABLE business_hold_transaction ALTER COLUMN amount_dep SET DEFAULT 0;
ALTER TABLE business_hold_transaction RENAME COLUMN amount_fd TO amount;
ALTER TABLE business_hold_transaction ALTER COLUMN amount DROP DEFAULT;

/* Account daily balance */
ALTER TABLE business_account_daily_balance RENAME COLUMN posted_balance TO posted_balance_dep;
ALTER TABLE business_account_daily_balance ALTER COLUMN posted_balance_dep SET DEFAULT 0;
ALTER TABLE business_account_daily_balance RENAME COLUMN posted_balance_fd TO posted_balance;
ALTER TABLE business_account_daily_balance ALTER COLUMN posted_balance DROP DEFAULT;

ALTER TABLE business_account_daily_balance RENAME COLUMN amount_credited TO amount_credited_dep;
ALTER TABLE business_account_daily_balance ALTER COLUMN amount_credited_dep SET DEFAULT 0;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_credited_fd TO amount_credited;
ALTER TABLE business_account_daily_balance ALTER COLUMN amount_credited DROP DEFAULT;

ALTER TABLE business_account_daily_balance RENAME COLUMN amount_debited TO amount_debited_dep;
ALTER TABLE business_account_daily_balance ALTER COLUMN amount_debited_dep SET DEFAULT 0;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_debited_fd TO amount_debited;
ALTER TABLE business_account_daily_balance ALTER COLUMN amount_debited DROP DEFAULT;

/* Daily transaction stats */
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_requested TO amount_requested_dep;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_requested_dep SET DEFAULT 0;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_requested_fd TO amount_requested;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_requested DROP DEFAULT;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_paid TO amount_paid_dep;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_paid_dep SET DEFAULT 0;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_paid_fd TO amount_paid;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_paid DROP DEFAULT;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_sent TO amount_sent_dep;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_sent_dep SET DEFAULT 0;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_sent_fd TO amount_sent;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_sent DROP DEFAULT;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_credited TO amount_credited_dep;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_credited_dep SET DEFAULT 0;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_credited_fd TO amount_credited;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_credited DROP DEFAULT;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_debited TO amount_debited_dep;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_debited_dep SET DEFAULT 0;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_debited_fd TO amount_debited;
ALTER TABLE business_daily_transaction_stats ALTER COLUMN amount_debited DROP DEFAULT;

/* Alloy results */
ALTER TABLE business_transaction_alloy_result RENAME COLUMN amount TO amount_dep;
ALTER TABLE business_transaction_alloy_result ALTER COLUMN amount_dep SET DEFAULT 0;
ALTER TABLE business_transaction_alloy_result RENAME COLUMN amount_fd TO amount;
ALTER TABLE business_transaction_alloy_result ALTER COLUMN amount DROP DEFAULT;

-- +goose Down
ALTER TABLE business_transaction_alloy_result RENAME COLUMN amount TO amount_fd;
ALTER TABLE business_transaction_alloy_result RENAME COLUMN amount_dep TO amount;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_debited TO amount_debited_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_debited_dep TO amount_debited;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_credited TO amount_credited_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_credited_dep TO amount_credited;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_sent TO amount_sent_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_sent_dep TO amount_sent;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_paid TO amount_paid_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_paid_dep TO amount_paid;

ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_requested TO amount_requested_fd;
ALTER TABLE business_daily_transaction_stats RENAME COLUMN amount_requested_dep TO amount_requested;

ALTER TABLE business_account_daily_balance RENAME COLUMN posted_balance TO posted_balance_fd;
ALTER TABLE business_account_daily_balance RENAME COLUMN posted_balance_dep TO posted_balance;

ALTER TABLE business_account_daily_balance RENAME COLUMN amount_credited TO amount_credited_fd;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_credited_dep TO amount_credited;

ALTER TABLE business_account_daily_balance RENAME COLUMN amount_debited TO amount_debited_fd;
ALTER TABLE business_account_daily_balance RENAME COLUMN amount_debited_dep TO amount_debited;

ALTER TABLE business_hold_transaction RENAME COLUMN amount TO amount_fd;
ALTER TABLE business_hold_transaction RENAME COLUMN amount_dep TO amount;

ALTER TABLE business_card_transaction RENAME COLUMN local_amount TO local_amount_fd;
ALTER TABLE business_card_transaction RENAME COLUMN local_amount_dep TO local_amount;

ALTER TABLE business_card_transaction RENAME COLUMN auth_amount TO auth_amount_fd;
ALTER TABLE business_card_transaction RENAME COLUMN auth_amount_dep TO auth_amount;

ALTER TABLE business_transaction RENAME COLUMN amount TO amount_fd;
ALTER TABLE business_transaction RENAME COLUMN amount_dep TO amount;
