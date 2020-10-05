-- +goose Up
/* Business request transfer table */
CREATE TABLE business_money_request
  ( 
    id                          UUID PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    created_user_id             UUID NOT NULL REFERENCES wise_user (id), 
    business_id                 UUID NOT NULL REFERENCES business (id), 
    contact_id                  UUID REFERENCES business_contact (id) DEFAULT NULL, 
    amount                      DOUBLE PRECISION NOT NULL, 
    currency                    TEXT NOT NULL, 
    notes                       TEXT DEFAULT NULL, 
    request_status              TEXT NOT NULL DEFAULT 'pending',
    message_id                  TEXT NOT NULL, 
    created                     TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp, 
    modified                    TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp 
  ); 

CREATE INDEX business_money_request_created_user_id_fkey ON business_money_request (created_user_id);
CREATE INDEX business_money_request_business_id_fkey ON business_money_request (business_id);
CREATE INDEX business_money_request_contact_id_fkey ON business_money_request (contact_id);

CREATE TRIGGER update_business_money_request_modified BEFORE UPDATE ON business_money_request FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_money_request_modified on business_money_request;
DROP INDEX business_money_request_created_user_id_fkey;
DROP INDEX business_money_request_business_id_fkey;
DROP INDEX business_money_request_contact_id_fkey;
DROP TABLE business_money_request;
