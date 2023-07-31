-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `tr_sample_reserve` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `order_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `office_id` INT DEFAULT NULL,
    `order_item_id` INT DEFAULT NULL,
    `qty` INT DEFAULT NULL,
    `note` VARCHAR(255) DEFAULT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` VARCHAR(255) DEFAULT NULL,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_by` VARCHAR(255) DEFAULT NULL,
    `deleted_at` TIMESTAMP DEFAULT NULL,
    `deleted_by` VARCHAR(255) DEFAULT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `tr_sample_reserve`;
-- +goose StatementEnd
