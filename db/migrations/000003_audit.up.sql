--
-- A relation that hold audit records that can help us piece
-- together exactly what happened in a request if necessary
-- after the fact. It can also, for example, be used to
-- drive internal security programs tasked with looking for
-- suspicious activity.
--
CREATE TABLE audit_records (
    id                 BIGSERIAL       PRIMARY KEY,

    -- action taken, for example "created"
    action             TEXT            NOT NULL
        CHECK (char_length(action) <= 50),

    created_at         TIMESTAMPTZ     NOT NULL DEFAULT now(),
    data               JSONB           NOT NULL,
    origin_ip          CIDR            NOT NULL,

    -- resource ID and type, for example "ride" ID 123
    resource_id        BIGINT          NOT NULL,
    resource_type      TEXT            NOT NULL
        CHECK (char_length(resource_type) <= 50),

    user_id            BIGINT          NOT NULL
        REFERENCES users ON DELETE RESTRICT
);