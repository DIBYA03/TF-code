-- +goose Up
ALTER TABLE wise_user ADD COLUMN subscription_status TEXT DEFAULT 'pending_acceptance';
ALTER TABLE business ADD COLUMN subscription_decision_date TIMESTAMP WITH time zone DEFAULT NULL;
ALTER TABLE business ADD COLUMN subscription_status TEXT DEFAULT 'pending_acceptance';
ALTER TABLE business ADD COLUMN subscription_start_date DATE DEFAULT NULL;

-- +goose Down
ALTER TABLE business DROP COLUMN subscription_start_date;
ALTER TABLE business DROP COLUMN subscription_status;
ALTER TABLE business DROP COLUMN subscription_decision_date;
ALTER TABLE wise_user DROP COLUMN subscription_status;