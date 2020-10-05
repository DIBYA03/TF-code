-- +goose Up
ALTER TABLE point_of_sale RENAME TO card_reader;

-- +goose Down
ALTER TABLE card_reader RENAME TO point_of_sale;
