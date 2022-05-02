CREATE TABLE IF NOT EXISTS "pairs"
(
    "pool_address"   TEXT             NOT NULL,
    "base_address"   TEXT             NOT NULL,
    "quote_address"  TEXT             NOT NULL,
    "is_base_token0" BOOLEAN          NOT NULL,
    "price"          DOUBLE PRECISION NOT NULL,
    "24h_change"     DOUBLE PRECISION NOT NULL,
    "24h_low"        DOUBLE PRECISION NOT NULL,
    "24h_high"       DOUBLE PRECISION NOT NULL,
    "24h_volume"     DOUBLE PRECISION NOT NULL,
    PRIMARY KEY ("pool_address")
);

CREATE TABLE IF NOT EXISTS "pair_subscriptions"
(
    "id"           SERIAL,
    "server_id"    TEXT NOT NULL,
    "pool_address" TEXT NOT NULL,
    PRIMARY KEY ("pool_address")
);

ALTER TABLE pair_subscriptions ADD CONSTRAINT pool_address_FK FOREIGN KEY ("pool_address") REFERENCES pairs("pool_address");