-- +goose Up
CREATE TABLE external_bank_account_owner ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    external_bank_account_id   UUID NOT NULL,
    account_holder_name        jsonb NOT NULL DEFAULT '[]'::jsonb,
    phone                      jsonb DEFAULT NULL,
    email                      jsonb DEFAULT NULL,
    owner_address              jsonb DEFAULT NULL,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX external_bank_account_owner_account_id_idx ON external_bank_account_owner(external_bank_account_id);

-- +goose Down
DROP INDEX external_bank_account_owner_account_id_idx;
DROP TABLE external_bank_account_owner;