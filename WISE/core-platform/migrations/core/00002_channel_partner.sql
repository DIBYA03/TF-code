-- +goose Up
/* Partner table */
CREATE TABLE channel_partner (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	channel_name                    text NOT NULL,
	code                            text NOT NULL,
	granted_license_count           integer DEFAULT 0,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_channel_partner_modified BEFORE UPDATE ON channel_partner FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_channel_partner_modified on channel_partner;
DROP TABLE channel_partner;
