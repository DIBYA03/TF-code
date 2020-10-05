-- +goose Up
/* Add pgcrypto extension  */
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION public.generate_uid(size integer) RETURNS text
    LANGUAGE plpgsql
    AS $$
 DECLARE
   characters TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   bytes BYTEA := gen_random_bytes(size);
   l INT := length(characters);
   i INT := 0;
   output TEXT := '';
 BEGIN
   WHILE i < size LOOP
     output := output || substr(characters, get_byte(bytes, i) % l + 1, 1);
     i := i + 1;
   END LOOP;
   RETURN output;
 END;
 $$;

-- +goose StatementEnd

/* Business request payment table */
CREATE TABLE business_money_request_payment
  ( 
    id                          UUID PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(), 
    request_id                  UUID NOT NULL REFERENCES business_money_request (id), 
    source_payment_id           TEXT NOT NULL, 
    status                      TEXT NOT NULL,
    token                       TEXT UNIQUE DEFAULT generate_uid(16), 
    expiration_date             timestamp with time zone NOT NULL,
    created                     TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp, 
    modified                    TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp 
  ); 

CREATE INDEX business_money_request_payment_request_id_fkey ON business_money_request_payment (request_id);

CREATE TRIGGER update_business_money_request_payment_modified BEFORE UPDATE ON business_money_request_payment FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_money_request_payment_modified on business_money_request_payment;
DROP INDEX business_money_request_payment_request_id_fkey;
DROP TABLE business_money_request_payment;
DROP FUNCTION public.generate_uid;
