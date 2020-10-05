-- +goose Up
/* Business table */
CREATE TABLE business (
	id 				uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	owner_id			uuid NOT NULL REFERENCES wise_user (id),
	employer_number			text NOT NULL,
	legal_name			text NOT NULL,
	dba				jsonb NOT NULL DEFAULT '[]'::jsonb,
	entity_type			text NOT NULL,
	industry_type			text NOT NULL,
	handles_cash			boolean NOT NULL DEFAULT false,
	tax_id				text NOT NULL,
	tax_id_type			text NOT NULL,
	origin_country			text DEFAULT NULL,
	origin_state			text DEFAULT NULL,
	origin_date			date DEFAULT NULL,
	kyc_status                      text NOT NULL DEFAULT 'notStarted',
	purpose				text NOT NULL,
	operation_type			text DEFAULT NULL,
	email				text DEFAULT NULL,
	email_verified   		text NOT NULL DEFAULT false,
	phone 				text DEFAULT NULL,
	phone_verified			boolean NOT NULL DEFAULT false,
	activity_type			jsonb NOT NULL DEFAULT '[]'::jsonb,
	legal_address                   jsonb DEFAULT NULL,
	headquarter_address             jsonb DEFAULT NULL,
	mailing_address                 jsonb DEFAULT NULL,
	is_restricted                   boolean NOT NULL DEFAULT false,
	formation_document_id           uuid DEFAULT NULL,
	deactivated			timestamp with time zone DEFAULT NULL,
	created				timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
	modified 			timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX business_owner_id_fkey ON business (owner_id);
CREATE UNIQUE INDEX business_tax_id_tax_id_type_idx ON business(tax_id, tax_id_type);

CREATE TRIGGER update_business_modified BEFORE UPDATE ON business FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_business_modified on business;
DROP INDEX business_owner_id_fkey;
DROP TABLE business;
