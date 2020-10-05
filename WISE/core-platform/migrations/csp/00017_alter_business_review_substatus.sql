-- +goose Up 
ALTER TABLE business ADD COLUMN review_substatus text NOT NULL DEFAULT 'wise';

-- +goose Down
ALTER TABLE business DROP COLUMN IF EXISTS review_substatus;
