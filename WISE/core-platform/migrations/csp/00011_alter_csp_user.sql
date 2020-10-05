-- +goose Up
ALTER TABLE csp_user ADD COLUMN cognito_id uuid NOT NULL;
ALTER TABLE csp_user ADD COLUMN picture TEXT DEFAULT NULL;
ALTER TABLE csp_user ADD COLUMN active BOOLEAN DEFAULT false;

ALTER TABLE csp_user ALTER COLUMN phone DROP NOT NULL;
ALTER TABLE csp_user ALTER COLUMN phone_verified DROP NOT NULL;
ALTER TABLE csp_user DROP CONSTRAINT csp_user_phone_key;

ALTER TABLE csp_user ALTER COLUMN email SET NOT NULL;

-- +goose Down
ALTER TABLE csp_user ALTER COLUMN email DROP NOT NULL;

ALTER TABLE csp_user ADD CONSTRAINT csp_user_phone_key UNIQUE (phone);
ALTER TABLE csp_user ALTER COLUMN phone_verified SET NOT NULL;
ALTER TABLE csp_user ALTER COLUMN phone SET NOT NULL;

ALTER TABLE csp_user DROP COLUMN active;
ALTER TABLE csp_user DROP COLUMN picture;
ALTER TABLE csp_user DROP COLUMN cognito_id;
