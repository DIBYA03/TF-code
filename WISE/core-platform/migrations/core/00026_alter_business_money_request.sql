-- +goose Up
/* PLATFORM-355: Add request type and pos ID to track terminal from which request is originating  */
ALTER TABLE business_money_request ADD COLUMN request_type TEXT DEFAULT NULL;
ALTER TABLE business_money_request ADD COLUMN pos_id uuid REFERENCES point_of_sale (id);

-- +goose Down
ALTER TABLE business_money_request DROP COLUMN request_type;
ALTER TABLE business_money_request DROP COLUMN pos_id;
