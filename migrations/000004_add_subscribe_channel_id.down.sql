BEGIN;

ALTER TABLE "pair_subscriptions" DROP COLUMN "channel_id";

COMMIT;
