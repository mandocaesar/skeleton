-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "public"."tr_sample" (
	"id" SERIAL PRIMARY KEY,
	"reference" VARCHAR DEFAULT NULL,
	"shipping_fee" FLOAT DEFAULT NULL,
	"insurance_fee" FLOAT NOT NULL,
	"adjustment_fee" FLOAT NOT NULL,
	"total_price" FLOAT NOT NULL,
	"total_price_old" FLOAT NOT NULL,
    "order_date" TIMESTAMPTZ NOT NULL default now(),
	"order_id" INT DEFAULT NULL,
    "note" VARCHAR DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL default now(),
	"created_by" VARCHAR DEFAULT NULL,
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"updated_by" VARCHAR DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ DEFAULT NULL,
	"deleted_by" VARCHAR DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS "public"."tr_sample_reserve" (
	"id" SERIAL PRIMARY KEY,
    "order_date" TIMESTAMPTZ NOT NULL default now(),
	"office_id" INT DEFAULT NULL,
    "order_item_id" INT DEFAULT NULL,
    "qty" INT DEFAULT NULL,
    "note" VARCHAR DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL default now(),
	"created_by" VARCHAR DEFAULT NULL,
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"updated_by" VARCHAR DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ DEFAULT NULL,
	"deleted_by" VARCHAR DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."tr_sample";
DROP TABLE IF EXISTS "public"."tr_sample_reserve";
-- +goose StatementEnd