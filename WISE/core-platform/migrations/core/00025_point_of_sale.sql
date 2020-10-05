-- +goose Up
CREATE TABLE point_of_sale ( 
    id                         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    business_id                UUID REFERENCES business (id), 
    alias                      TEXT DEFAULT NULL, 
    device_type                TEXT DEFAULT NULL, 
    serial_number              TEXT NOT NULL, 
    last_connected             TIMESTAMP WITH time zone DEFAULT NULL, 
    deactivated                TIMESTAMP WITH time zone DEFAULT NULL, 
    created                    TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    modified                   TIMESTAMP WITH time zone NOT NULL DEFAULT CURRENT_TIMESTAMP 
); 

CREATE INDEX point_of_sale_business_id_fk ON point_of_sale(business_id);
CREATE UNIQUE INDEX point_of_sale_serial_number_business_id_idx ON point_of_sale(serial_number, business_id);

-- +goose Down
DROP INDEX point_of_sale_serial_number_business_id_idx;
DROP INDEX point_of_sale_business_id_fk;
DROP TABLE point_of_sale;
