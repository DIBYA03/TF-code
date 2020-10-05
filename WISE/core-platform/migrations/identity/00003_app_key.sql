-- +goose Up
/* Wise user table */
CREATE TABLE app_key (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	app_id                          uuid NOT NULL REFERENCES app (id),
	key_name                        text NOT NULL,
	api_key                         text NOT NULL,
	api_secret                      text NOT NULL,
	description                     text  DEFAULT NULL,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX app_key_app_id_fkey ON app_key (app_id);
CREATE TRIGGER update_app_key_modified BEFORE UPDATE ON app_key FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_app_key_modified on app_key;
DROP INDEX app_key_app_id_fkey;
DROP TABLE app_key;
