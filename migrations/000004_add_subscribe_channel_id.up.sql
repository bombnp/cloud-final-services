BEGIN;

ALTER TABLE "pair_subscriptions" ADD COLUMN "channel_id" TEXT;

COMMIT;
