-- +goose Up
ALTER TABLE business_notes DROP CONSTRAINT business_notes_business_id_key;
ALTER TABLE business_notes DROP CONSTRAINT business_notes_user_id_key;

ALTER TABLE business ADD COLUMN accounting_provider text DEFAULT '';
ALTER TABLE business ADD COLUMN payroll_provider text DEFAULT '';
ALTER TABLE business ADD COLUMN payroll_type text DEFAULT ''; -- 1099,s 'w2' 'both'
ALTER TABLE business ADD COLUMN business_description text DEFAULT '';

-- +goose Down
ALTER TABLE business_notes ADD CONSTRAINT business_notes_user_id_key UNIQUE (user_id);
ALTER TABLE business_notes ADD CONSTRAINT business_notes_business_id_key UNIQUE (business_id);

ALTER TABLE business DROP COLUMN IF EXISTS accounting_provider;
ALTER TABLE business DROP COLUMN IF EXISTS payroll_provider; 
ALTER TABLE business DROP COLUMN IF EXISTS payroll_type;  
ALTER TABLE business DROP COLUMN IF EXISTS business_description;
