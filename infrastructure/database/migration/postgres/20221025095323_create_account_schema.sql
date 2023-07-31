-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."ms_account" (
	"id" SERIAL PRIMARY KEY,
	"customer_id" INT NOT NULL,
	"username" VARCHAR NOT NULL,
	"hash" VARCHAR NOT NULL,
	"is_active" BOOLEAN NOT NULL DEFAULT FALSE,
	"created_at" TIMESTAMPTZ NOT NULL default now(),
	"created_by" VARCHAR DEFAULT NULL,
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"updated_by" VARCHAR DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ DEFAULT NULL,
	"deleted_by" VARCHAR DEFAULT NULL
);

CREATE UNIQUE INDEX "ms_account_username_unique_idx" ON "ms_account" ("username") WHERE deleted_at IS NULL;
CREATE INDEX "ms_account_customer_id_idx" ON "ms_account" ("customer_id");

CREATE TABLE IF NOT EXISTS "public"."ms_customer" (
	"id" SERIAL PRIMARY KEY,
	"account_id" INT,
	"name" VARCHAR NOT NULL,
	"brand_name" VARCHAR NOT NULL,
	"email" VARCHAR NOT NULL,
	"phone" VARCHAR NOT NULL,
	"code" VARCHAR NOT NULL,
	"address" VARCHAR NOT NULL,
	"zip" INT DEFAULT NULL,
	"province_id" INT DEFAULT NULL,
	"district_id" INT DEFAULT NULL,
	"sub_district_id" INT DEFAULT NULL,
	"is_order_bypass" BOOLEAN NOT NULL DEFAULT FALSE,
	"created_at" TIMESTAMPTZ NOT NULL default now(),
	"created_by" VARCHAR DEFAULT NULL,
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"updated_by" VARCHAR DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ DEFAULT NULL,
	"deleted_by" VARCHAR DEFAULT NULL
);

CREATE INDEX "ms_customer_account_id_idx" ON "ms_customer" ("account_id");

CREATE TABLE IF NOT EXISTS "public"."ms_address" (
	"id" SERIAL PRIMARY KEY,
	"customer_id" INT,
	"name" VARCHAR NOT NULL,
	"phone" VARCHAR NOT NULL,
	"email" VARCHAR,
	"address" VARCHAR NOT NULL,
	"city" VARCHAR NULL,
	"zip" INT DEFAULT NULL,
	"province_id" INT DEFAULT NULL,
	"district_id" INT DEFAULT NULL,
	"sub_district_id" INT DEFAULT NULL,
	"latitude" FLOAT DEFAULT NULL,
	"longitude" FLOAT DEFAULT NULL,
	"label" VARCHAR,
	"is_default" BOOLEAN NOT NULL DEFAULT FALSE,
	"created_at" TIMESTAMPTZ NOT NULL default now(),
	"created_by" VARCHAR DEFAULT NULL,
	"updated_at" TIMESTAMPTZ DEFAULT now(),
	"updated_by" VARCHAR DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ DEFAULT NULL,
	"deleted_by" VARCHAR DEFAULT NULL
);

CREATE INDEX "ms_address_customer_id_idx" ON "ms_address" ("customer_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."ms_account";
DROP TABLE IF EXISTS "public"."ms_customer";
DROP TABLE IF EXISTS "public"."ms_address";
-- +goose StatementEnd
