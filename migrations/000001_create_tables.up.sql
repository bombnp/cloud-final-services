BEGIN;

CREATE TABLE IF NOT EXISTS "tokens"
(
    "address" TEXT NOT NULL,
    "symbol"  TEXT NOT NULL,
    "icon"    TEXT,
    PRIMARY KEY ("address")
);

CREATE TABLE IF NOT EXISTS "pairs"
(
    "pool_address"   TEXT    NOT NULL,
    "base_address"   TEXT    NOT NULL,
    "quote_address"  TEXT    NOT NULL,
    "is_base_token0" BOOLEAN NOT NULL,
    PRIMARY KEY ("pool_address"),
    CONSTRAINT fk_base_token FOREIGN KEY (base_address) REFERENCES tokens (address) ON DELETE RESTRICT,
    CONSTRAINT fk_quote_token FOREIGN KEY (quote_address) REFERENCES tokens (address) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS "pair_subscriptions"
(
    "id"           SERIAL,
    "server_id"    TEXT NOT NULL,
    "pool_address" TEXT NOT NULL,
    "type"         TEXT NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT pool_address_FK FOREIGN KEY ("pool_address") REFERENCES pairs ("pool_address") ON DELETE CASCADE
);

COMMIT;
