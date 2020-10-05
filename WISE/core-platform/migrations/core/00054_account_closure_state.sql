-- +goose Up
CREATE TABLE account_closure_state ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    account_closure_request_id UUID REFERENCES account_closure_request (id),
    closure_state              TEXT NOT NULL,
    item_id                    TEXT,
    description                TEXT,      
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

/*
 * Index account_closure_request_id
 */
CREATE INDEX account_closure_state_closure_request_id_fk ON account_closure_state (account_closure_request_id);

-- +goose Down
DROP INDEX account_closure_state_closure_request_id_fk

DROP TABLE account_closure_state;
