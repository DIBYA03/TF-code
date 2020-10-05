-- +goose Up
/* Email table */
CREATE TABLE email (
    id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    email_address               text NOT NULL,
    email_status                text NOT NULL,
    email_type                  text NOT NULL,
    created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX email_address_status_idx ON email(email_address, email_status);

CREATE TRIGGER update_email_modified BEFORE UPDATE ON email FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

ALTER TABLE business ADD COLUMN email_id uuid REFERENCES email(id) DEFAULT NULL;
ALTER TABLE consumer ADD COLUMN email_id uuid REFERENCES email(id) DEFAULT NULL;
ALTER TABLE business_contact ADD COLUMN email_id uuid REFERENCES email(id) DEFAULT NULL;

-- +goose Down
ALTER TABLE consumer DROP COLUMN email_id;
ALTER TABLE business DROP COLUMN email_id;
ALTER TABLE business_contact DROP COLUMN email_id;

DROP TRIGGER IF EXISTS update_email_modified on email;

DROP INDEX email_address_status_idx;
DROP TABLE email;
