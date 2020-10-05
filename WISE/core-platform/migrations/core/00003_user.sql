-- +goose Up
/* Wise user table */
CREATE TABLE wise_user (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                     uuid UNIQUE NOT NULL REFERENCES consumer (id),
	identity_id                     uuid UNIQUE NOT NULL,
	partner_id                      uuid REFERENCES channel_partner (id) DEFAULT NULL,
	email                           text DEFAULT NULL,
	email_verified                  boolean NOT NULL DEFAULT false,
	phone                           text UNIQUE NOT NULL,
	phone_verified                  boolean NOT NULL DEFAULT false,
	notification                    jsonb NOT NULL DEFAULT
	                                '{'
	                                        '"transfers": true,'
	                                        '"transactions": true,'
	                                        '"contacts": true'
	                                '}'::jsonb,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX wise_user_consumer_id_fkey ON wise_user (consumer_id);
CREATE INDEX wise_user_partner_id_fkey ON wise_user (partner_id);
CREATE INDEX wise_user_identity_id_idx ON wise_user (identity_id);

CREATE TRIGGER update_wise_user_modified BEFORE UPDATE ON wise_user FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_wise_user_modified on wise_user;
DROP INDEX wise_user_identity_id_idx;
DROP INDEX wise_user_partner_id_fkey;
DROP INDEX wise_user_consumer_id_fkey;
DROP TABLE wise_user;
