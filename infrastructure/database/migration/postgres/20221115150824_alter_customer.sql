-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS "ms_customer_account_id_idx";

ALTER TABLE "public"."ms_customer" DROP COLUMN IF EXISTS "account_id";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."ms_customer" ADD COLUMN IF NOT EXISTS "account_id" INT;

CREATE INDEX "ms_customer_account_id_idx" ON "ms_customer" ("account_id");
-- +goose StatementEnd
