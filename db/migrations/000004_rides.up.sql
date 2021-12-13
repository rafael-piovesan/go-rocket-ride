--
-- A relation representing a single ride by a user.
-- Notably, it holds the ID of a successful charge to
-- Stripe after we have one.
--
CREATE TABLE rides (
    id                 BIGSERIAL       PRIMARY KEY,
    created_at         TIMESTAMPTZ     NOT NULL DEFAULT now(),

    -- Store a reference to the idempotency key so that we can recover an
    -- already-created ride. Note that idempotency keys are not stored
    -- permanently, so make sure to SET NULL when a referenced key is being
    -- reaped.
    idempotency_key_id BIGINT          UNIQUE
        REFERENCES idempotency_keys ON DELETE SET NULL,

    -- origin and destination latitudes and longitudes
    origin_lat       NUMERIC(13, 10) NOT NULL,
    origin_lon       NUMERIC(13, 10) NOT NULL,
    target_lat       NUMERIC(13, 10) NOT NULL,
    target_lon       NUMERIC(13, 10) NOT NULL,

    -- ID of Stripe charge like ch_123; NULL until we have one
    stripe_charge_id   TEXT            UNIQUE
        CHECK (char_length(stripe_charge_id) <= 50),

    user_id            BIGINT          NOT NULL
        REFERENCES users ON DELETE RESTRICT
);

CREATE INDEX rides_idempotency_key_id
    ON rides (idempotency_key_id)
    WHERE idempotency_key_id IS NOT NULL;