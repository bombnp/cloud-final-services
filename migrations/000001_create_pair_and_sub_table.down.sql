BEGIN;

ALTER TABLE pair_subscriptions DROP CONSTRAINT IF EXISTS pool_address_FK;

DROP TABLE IF EXISTS pairs;
DROP TABLE IF EXISTS pair_subscriptions;

COMMIT;
