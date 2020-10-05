-- +goose Up
CREATE TABLE external_account_verification_result ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID NOT NULL,
    external_bank_account_id   UUID NOT NULL,
    source_ip_address          TEXT NOT NULL,
    access_token               TEXT NOT NULL,
    partner_item_id            TEXT NOT NULL,
    verification_status        TEXT NOT NULL,
    verification_result        jsonb NOT NULL DEFAULT '[]'::jsonb,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
); 

CREATE INDEX external_account_verification_result_business_id_idx ON external_account_verification_result(business_id);

-- +goose Down
DROP INDEX external_account_verification_result_business_id_idx;
DROP TABLE external_account_verification_result;