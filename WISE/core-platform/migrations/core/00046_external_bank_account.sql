-- +goose Up
CREATE TABLE external_bank_account ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID NOT NULL,
    linked_account_id          UUID DEFAULT NULL,
    partner_account_id         TEXT NOT NULL,
    partner_name               TEXT NOT NULL,
    account_name               TEXT NOT NULL,
    official_account_name      TEXT NOT NULL,
    account_type               TEXT NOT NULL,
    account_subtype            TEXT NOT NULL,
    account_number             TEXT NOT NULL,
    routing_number             TEXT NOT NULL,
    wire_routing               TEXT NOT NULL,
    available_balance          double precision NOT NULL,
    posted_balance             double precision NOT NULL,
    currency                   text NOT NULL,
    last_login                 timestamp with time zone DEFAULT NULL,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX external_bank_account_business_id_idx ON external_bank_account(business_id);

-- +goose Down
DROP INDEX external_bank_account_business_id_idx;
DROP TABLE external_bank_account;