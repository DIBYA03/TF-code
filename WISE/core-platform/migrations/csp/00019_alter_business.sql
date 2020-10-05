-- +goose Up 
ALTER TABLE business ADD COLUMN subscribed_agent_id UUID DEFAULT NULL REFERENCES csp_user(id);

-- +goose Down
ALTER TABLE business DROP COLUMN IF EXISTS subscribed_agent_id;
