-- +goose Up
/* Business member table  */
CREATE TABLE business_member  ( 
	id                              uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	consumer_id                     uuid NOT NULL REFERENCES consumer (id),
	business_id                     uuid NOT NULL REFERENCES business (id),
	title_type                      text NOT NULL,
	title_other                     text DEFAULT NULL,
	ownership                       integer NOT NULL DEFAULT 0,
	is_controlling_manager          boolean NOT NULL DEFAULT false,
	deactivated                     timestamp with time zone DEFAULT NULL,
	created                         timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified                        timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_member_business_id_fkey ON business_member (business_id);
CREATE INDEX business_member_consumer_id_fkey ON business_member (consumer_id);
CREATE UNIQUE INDEX business_member_business_id_consumer_id_idx ON business_member (business_id, consumer_id);

CREATE TRIGGER update_business_member_modified BEFORE UPDATE ON business_member FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_member_modified on business_member;
DROP INDEX business_member_business_id_consumer_id_idx;
DROP INDEX business_member_consumer_id_fkey;
DROP INDEX business_member_business_id_fkey;
DROP TABLE business_member;
