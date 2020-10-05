-- +goose Up
CREATE TABLE user_activity ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    entity_id                  UUID NOT NULL REFERENCES wise_user (id), 
    activity_type              TEXT NOT NULL, 
    activity_action            TEXT DEFAULT NULL, 
    resource_id                TEXT DEFAULT NULL, 
    metadata                   JSONB NOT NULL DEFAULT '{}'::jsonb, 
    activity_date              TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp, 
    created                    TIMESTAMP WITH TIME zone NOT NULL DEFAULT current_timestamp 
  ); 

CREATE INDEX user_activity_entity_id_idx ON user_activity (entity_id);
CREATE INDEX user_activity_entity_id_activity_type_idx ON user_activity (entity_id, activity_type);
CREATE INDEX user_activity_type_resource_id_idx ON user_activity (activity_type, resource_id);

-- +goose Down
DROP INDEX user_activity_entity_id_idx;
DROP INDEX user_activity_entity_id_activity_type_idx;
DROP INDEX user_activity_type_resource_id_idx;
DROP TABLE user_activity;
