-- +goose Up
CREATE TABLE business_card_reissue_history ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID REFERENCES business (id),
    card_id                    UUID REFERENCES business_bank_card (id),
    reason                     TEXT NOT NULL,
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX business_card_reissue_history_business_id_fk ON business_card_reissue_history(business_id);
CREATE INDEX business_card_reissue_history_card_id_fk ON business_card_reissue_history(card_id);

-- +goose Down
DROP INDEX business_card_reissue_history_card_id_fk;
DROP INDEX business_card_reissue_history_business_id_fk;
DROP TABLE business_card_reissue_history;