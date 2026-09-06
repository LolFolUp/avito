-- +goose Up
DROP INDEX menu_items_restaurant_available_idx;

ALTER TABLE menu_items
    DROP CONSTRAINT menu_items_pkey,
    ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid(),
    ADD COLUMN version BIGINT NOT NULL DEFAULT 1,
    DROP COLUMN external_id,
    ADD CONSTRAINT menu_items_pkey PRIMARY KEY (id),
    ADD CONSTRAINT menu_items_restaurant_id_id_unique UNIQUE (restaurant_id, id),
    ADD CONSTRAINT menu_items_version_positive CHECK (version > 0);

CREATE INDEX menu_items_restaurant_available_idx
    ON menu_items (restaurant_id, id)
    WHERE is_active = TRUE AND is_available = TRUE;

-- +goose Down
DROP INDEX menu_items_restaurant_available_idx;

ALTER TABLE menu_items
    DROP CONSTRAINT menu_items_version_positive,
    DROP CONSTRAINT menu_items_restaurant_id_id_unique,
    DROP CONSTRAINT menu_items_pkey,
    ADD COLUMN external_id VARCHAR(200),
    DROP COLUMN version,
    DROP COLUMN id;

UPDATE menu_items
SET external_id = gen_random_uuid()::TEXT;

ALTER TABLE menu_items
    ALTER COLUMN external_id SET NOT NULL,
    ADD CONSTRAINT menu_items_pkey PRIMARY KEY (restaurant_id, external_id);

CREATE INDEX menu_items_restaurant_available_idx
    ON menu_items (restaurant_id, external_id)
    WHERE is_active = TRUE AND is_available = TRUE;
