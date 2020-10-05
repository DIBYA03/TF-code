-- +goose Up
CREATE TABLE money_transfer_request ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    created_user_id            uuid NOT NULL REFERENCES wise_user (id),
    business_id                UUID REFERENCES business (id), 
    contact_id                 UUID REFERENCES business_contact (id) DEFAULT NULL, 
    money_transfer_id          UUID REFERENCES business_money_transfer (id) DEFAULT NULL, 
    request_mode               TEXT NOT NULL, 
    payment_token              TEXT UNIQUE DEFAULT generate_uid(16), 
    amount                     DECIMAL(19,4) NOT NULL,
    notes                      TEXT DEFAULT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX money_transfer_request_created_user_id_fkey ON money_transfer_request(created_user_id);
CREATE INDEX money_transfer_request_business_id_fkey ON money_transfer_request(business_id);
CREATE INDEX money_transfer_request_contact_id_fkey ON money_transfer_request(contact_id);
CREATE INDEX money_transfer_request_money_transfer_id_fkey ON money_transfer_request(money_transfer_id);

-- +goose Down
DROP INDEX money_transfer_request_money_transfer_id_fkey;
DROP INDEX money_transfer_request_contact_id_fkey;
DROP INDEX money_transfer_request_business_id_fkey;
DROP INDEX money_transfer_request_created_user_id_fkey;
DROP TABLE money_transfer_request;
