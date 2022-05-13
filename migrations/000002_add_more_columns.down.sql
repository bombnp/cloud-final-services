BEGIN;

ALTER TABLE "pairs" DROP COLUMN "dex";
ALTER TABLE "tokens" DROP COLUMN "name";

COMMIT;
