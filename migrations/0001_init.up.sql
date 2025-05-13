CREATE TABLE clients
(
    id           TEXT PRIMARY KEY,
    capacity     INTEGER     NOT NULL,
    rate_per_sec INTEGER     NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL
);
