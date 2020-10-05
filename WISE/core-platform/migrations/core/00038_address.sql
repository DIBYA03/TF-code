-- +goose Up
/* Address table */
CREATE TABLE address (
    id                          uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    consumer_id                 uuid REFERENCES consumer(id),
    contact_id                  uuid REFERENCES business_contact(id),
    business_id                 uuid REFERENCES business(id),
    street                      text NOT NULL,
    line2                       text NOT NULL DEFAULT '',
    locality                    text NOT NULL, --City, municipality
    admin_area                  text NOT NULL, --State, province, etc
    country                     text NOT NULL,
    postal_code                 text NOT NULL,
    latitude                    decimal(9,6) NOT NULL DEFAULT 0,
    longitude                   decimal(9,6) NOT NULL DEFAULT 0,
    address_type                text NOT NULL,
    address_state               text NOT NULL,
    created                     timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified                    timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP

    /* make sure one of these fields is set */
    CONSTRAINT has_one_fkey check(
        (
            (consumer_id IS NOT NULL)::integer +
            (contact_id IS NOT NULL)::integer +
            (business_id IS NOT NULL)::integer
        ) = 1
    )
);

CREATE INDEX address_consumer_id_idx ON address(consumer_id);
CREATE INDEX address_contact_id_idx ON address(contact_id);
CREATE INDEX address_business_id_idx ON address(business_id);

ALTER TABLE business_linked_payee ADD CONSTRAINT address_id_fkey FOREIGN KEY (address_id) REFERENCES address(id);

CREATE TRIGGER update_address_modified BEFORE UPDATE ON address FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- +goose Down
ALTER TABLE business_linked_payee DROP CONSTRAINT address_id_fkey;
DROP TRIGGER IF EXISTS update_address_modified on address;
DROP INDEX address_consumer_id_idx;
DROP INDEX address_contact_id_idx;
DROP INDEX address_business_id_idx;
DROP TABLE address;
