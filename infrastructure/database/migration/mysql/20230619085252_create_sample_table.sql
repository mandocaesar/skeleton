-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS `tr_sample` (
    `id` int AUTO_INCREMENT PRIMARY KEY,
    `reference` VARCHAR(255) DEFAULT NULL,
    `shipping_fee` FLOAT DEFAULT NULL,
    `insurance_fee` FLOAT NOT NULL,
    `adjustment_fee` FLOAT NOT NULL,
    `total_price` FLOAT NOT NULL,
    `total_price_old` FLOAT NOT NULL,
    `order_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `order_id` INT DEFAULT NULL,
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
DROP TABLE IF EXISTS `tr_sample`;


-- +goose StatementEnd
