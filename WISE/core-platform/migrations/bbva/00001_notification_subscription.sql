-- +goose Up
/* Add UUID extension  */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE notification_subscription (
	id                            uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	subscription_id               text NOT NULL,
	event_id                      text NOT NULL,
	event_type                    text NOT NULL,
	event_desc                    text NOT NULL,
	channel_type                  text NOT NULL,
	channel_url                   text NOT NULL,
	created                       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX notification_subscription_subscription_id_idx ON notification_subscription (subscription_id);
CREATE INDEX notification_subscription_event_id_idx ON notification_subscription (event_id);
CREATE INDEX notification_subscription_channel_url_event_type_idx ON notification_subscription (channel_url, event_type);

-- +goose Down
DROP INDEX notification_subscription_channel_url_event_type_idx;
DROP INDEX notification_subscription_event_id_idx;
DROP INDEX notification_subscription_subscription_id_idx;
DROP TABLE notification_subscription;
