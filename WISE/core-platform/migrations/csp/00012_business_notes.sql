-- +goose Up 
CREATE TABLE business_notes (
    id  uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid UNIQUE NOT NULL REFERENCES csp_user (id),
    business_id uuid UNIQUE  NOT NULL,
    notes text NOT NULL DEFAULT '',
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    modified timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER update_business_notes_modified BEFORE UPDATE ON business_notes FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE INDEX business_notes_csp_user_id_fkey ON business_notes(user_id);
-- +goose Down
DROP TRIGGER IF EXISTS update_business_notes_modified  on business_notes;
DROP INDEX business_notes_csp_user_id_fkey;
DROP TABLE business_notes;