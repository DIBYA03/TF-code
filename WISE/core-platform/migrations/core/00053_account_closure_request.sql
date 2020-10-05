-- +goose Up
CREATE TABLE account_closure_request ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID REFERENCES business (id),
    reason                     TEXT NOT NULL,
    status                     TEXT NOT NULL DEFAULT 'pending',
    csp_agent_id               uuid DEFAULT NULL,
    refund_amount              double precision,
    digital_check_number       TEXT,
    description                TEXT,
    closed                     TIMESTAMP WITH time zone, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

/*
 * Index business id
 */
CREATE INDEX account_closure_request_business_id_fk ON account_closure_request (business_id);

-- +goose Down
DROP INDEX account_closure_request_business_id_fk;

DROP TABLE account_closure_request;


