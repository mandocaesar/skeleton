-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."ms_region_province" (
	"id" INT PRIMARY KEY,
	"name" VARCHAR NOT NULL,
	"code" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."ms_region_district" (
	"id" INT PRIMARY KEY,
	"province_id" INT NOT NULL,
	"name" VARCHAR NOT NULL,
	"latitude" FLOAT DEFAULT NULL,
	"longitude" FLOAT DEFAULT NULL,
	"level" VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS "public"."ms_region_sub_district" (
	"id" INT PRIMARY KEY,
	"province_id" INT NOT NULL,
	"district_id" INT NOT NULL,
	"name" VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."ms_region_province";
DROP TABLE IF EXISTS "public"."ms_region_district";
DROP TABLE IF EXISTS "public"."ms_region_sub_district";
-- +goose StatementEnd
