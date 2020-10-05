-- +goose Up 
CREATE TABLE business (
    id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    business_id uuid UNIQUE  NOT NULL,
    business_name text NOT NULL,
    entity_type text NOT NULL,
    -- use to check in what status the business is e.g( account created | card created)
    process_status text NOT NULL DEFAULT 'created',
    review_status text NOT NULL DEFAULT 'pendingReview',
    idvs jsonb DEFAULT '[]'::jsonb,
    notes text NOT NULL DEFAULT '',
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_business_modified BEFORE UPDATE ON business FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE INDEX business_status_business_id_idx ON business(business_id);
CREATE INDEX business_status_status_idx ON business(review_status);

-- +goose Down
DROP TRIGGER IF EXISTS update_business_modified  on business;
DROP INDEX business_status_business_id_idx;
DROP INDEX business_status_status_idx;
DROP TABLE business;
