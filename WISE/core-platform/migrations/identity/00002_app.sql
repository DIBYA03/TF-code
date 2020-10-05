-- +goose Up
/* Wise user table */
CREATE TABLE app (
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	app_name                        text NOT NULL,
	description                     text  DEFAULT NULL,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_app_modified BEFORE UPDATE ON app FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_app_modified on app;
DROP TABLE app;
