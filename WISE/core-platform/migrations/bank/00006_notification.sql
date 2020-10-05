-- +goose Up
CREATE TABLE notification (
	id                            uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	entity_id                     uuid NOT NULL,
	entity_type                   text NOT NULL,
	bank_name                     text NOT NULL,
	source_id                     text NOT NULL,
	notification_type             text NOT NULL,
	notification_action           text NOT NULL,
	notification_attribute        text DEFAULT NULL,
	notification_version          text NOT NULL,
	send_counter                  int DEFAULT 0,
	created                       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	notification_data             jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX notification_entity_id_entity_type_idx ON notification (entity_id, entity_type);
CREATE UNIQUE INDEX notification_source_id_bank_name_idx ON notification (source_id, bank_name);

-- +goose Down
DROP INDEX notification_source_id_bank_name_idx;
DROP INDEX notification_entity_id_entity_type_idx;
DROP TABLE notification;
