-- +goose Up 
CREATE TABLE consumer_state (
    id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    consumer_id uuid  NOT NULL REFERENCES consumer(id),
    process_status text NOT NULL DEFAULT '',
    review_status text NOT NULL DEFAULT '',
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX consumer_state_consumer_id_fk ON consumer_state(consumer_id);

-- +goose Down
DROP INDEX consumer_state_consumer_id_fk;
DROP TABLE consumer_state;
