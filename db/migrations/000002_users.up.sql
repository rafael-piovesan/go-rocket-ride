--
-- A relation to hold records for every user of our app.
--
CREATE TABLE users (
    id                 BIGSERIAL       PRIMARY KEY,
    email              TEXT            NOT NULL UNIQUE
        CHECK (char_length(email) <= 255),

    -- Stripe customer record with an active credit card
    stripe_customer_id TEXT            NOT NULL UNIQUE
        CHECK (char_length(stripe_customer_id) <= 50)
);

CREATE INDEX users_email
    ON users (email);

--
-- Now that we have a users table, add a foreign key
-- constraint to idempotency_keys which we created above.
--
ALTER TABLE idempotency_keys
    ADD CONSTRAINT idempotency_keys_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;