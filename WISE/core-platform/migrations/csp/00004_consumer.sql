-- +goose Up 

CREATE TABLE consumer (
    id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    consumer_id uuid UNIQUE  NOT NULL,
    consumer_name text DEFAULT NULL,
    review_status text NOT NULL DEFAULT 'review',
    idvs jsonb DEFAULT '[]'::jsonb,
    notes text NOT NULL DEFAULT '',
    resolved timestamp with time zone DEFAULT NULL,
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_consumer_modified BEFORE UPDATE ON consumer FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE INDEX consumer_id_idx ON consumer(consumer_id);
CREATE INDEX consumer_review_status_idx ON consumer(review_status);

-- +goose Down
DROP TRIGGER IF EXISTS update_consumer_modified on consumer;
DROP INDEX consumer_id_idx;
DROP INDEX consumer_review_status_idx;
DROP TABLE consumer;
