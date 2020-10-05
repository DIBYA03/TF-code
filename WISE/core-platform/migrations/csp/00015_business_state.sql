
-- +goose Up 
CREATE TABLE business_state (
    id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    business_id uuid NOT NULL REFERENCES business (id),
    process_status text NOT NULL DEFAULT '',
    review_status text NOT NULL DEFAULT '',
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX   business_state_business_id_fk ON business_state(business_id);

-- +goose Down
DROP INDEX business_state_business_id_fk;
DROP TABLE business_state;

