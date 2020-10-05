-- +goose Up
CREATE TABLE business_contact (
	id                         uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	user_id                    uuid NOT NULL REFERENCES wise_user (id),
	business_id                uuid NOT NULL REFERENCES business (id),
	contact_category           text DEFAULT NULL,
	contact_type               text DEFAULT NULL,
	engagement                 text DEFAULT NULL,
	job_title                  text DEFAULT NULL,
	business_name              text DEFAULT NULL,
	first_name                 text DEFAULT NULL,
	last_name                  text DEFAULT NULL,
	phone_number               text DEFAULT NULL,
	email                      text DEFAULT NULL,
	mailing_address            jsonb DEFAULT NULL,
	deactivated                timestamp with time zone DEFAULT NULL,
	created                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified                   timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_contact_user_id_fkey ON business_contact (user_id);
CREATE INDEX business_contact_business_id_fkey ON business_contact (business_id);

CREATE TRIGGER update_business_contact_modified BEFORE UPDATE ON business_contact FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_contact_modified on business_contact;
DROP INDEX business_contact_user_id_fkey;
DROP INDEX business_contact_business_id_fkey;
DROP TABLE business_contact;
