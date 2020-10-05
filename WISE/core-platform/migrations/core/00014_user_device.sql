-- +goose Up
/* User Device table  */
CREATE TABLE user_device (
    id          UUID PRIMARY KEY NOT NULL DEFAULT Uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES wise_user (id),
    device_type TEXT NOT NULL,
    token_type  TEXT NOT NULL,
    token       TEXT NOT NULL,
    device_key  TEXT NOT NULL,
    language    TEXT NOT NULL DEFAULT 'en-US',
    created     TIMESTAMP WITH time zone DEFAULT CURRENT_TIMESTAMP,
    modified    TIMESTAMP WITH time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX user_device_user_id_fkey ON user_device(user_id);
CREATE INDEX user_device_token_idx ON user_device(token);
CREATE INDEX user_device_user_id_token_idx ON user_device(user_id, token);
CREATE INDEX user_device_device_key_idx ON user_device(device_key);

CREATE TRIGGER update_user_device_modified BEFORE UPDATE ON user_device FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_user_device_modified on user_device;
DROP INDEX user_device_user_id_fkey;
DROP INDEX user_device_token_idx;
DROP INDEX user_device_user_id_token_idx;
DROP INDEX user_device_device_key_idx;
DROP TABLE user_device;
