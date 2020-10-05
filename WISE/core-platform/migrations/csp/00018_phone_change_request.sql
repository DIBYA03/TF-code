-- +goose Up
CREATE TABLE phone_change_request ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    user_id                    UUID NOT NULL,
    old_phone                  TEXT NOT NULL,
    new_phone                  TEXT NOT NULL,
    originated_from            TEXT NOT NULL,
    csp_user_id                UUID DEFAULT NULL,
    verification_notes         TEXT NOT NULL,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
); 

CREATE INDEX phone_change_request_user_id_idx ON phone_change_request(user_id);

-- +goose Down
DROP INDEX phone_change_request_user_id_idx;
DROP TABLE phone_change_request;