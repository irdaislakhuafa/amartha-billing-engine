CREATE DATABASE IF NOT EXISTS `amartha_billing`;
USE `amartha_billing`;

-- loans
CREATE TABLE `loans` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `name` VARCHAR(255) NOT NULL,
 `description` TEXT NOT NULL,
 `interest_rate` DECIMAL(10, 2) NOT NULL COMMENT "per annum",
 `repayment_type` VARCHAR(255) NOT NULL COMMENT "weeks, months, years",
 `repayment_duration` INT NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

-- users
CREATE TABLE `users` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `name` VARCHAR(255) NOT NULL,
 `email` VARCHAR(255) NOT NULL,
 `password` VARCHAR(255) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

-- loan transactions
CREATE TABLE `loan_transactions` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `invoice_number` VARCHAR(255) NOT NULL,
 `notes` TEXT NOT NULL,
 `user_id` BIGINT NOT NULL COMMENT "refer to users.id",
 `user` JSON NOT NULL,
 `loan_id` BIGINT NOT NULL COMMENT "refer to loans.id",
 `loan` JSON NOT NULL,
 `amount` DECIMAL(10, 2) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

-- loans billing
CREATE TABLE `loans_billing` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `loan_transaction_id` BIGINT NOT NULL COMMENT "refer to loan_transactions.id",
 `bill_date` DATE NOT NULL,
 `principal_amount` DECIMAL(10, 2) NOT NULL,
 `principal_amount_paid` DECIMAL(10, 2) NOT NULL,
 `interest_amount` DECIMAL(10, 2) NOT NULL,
 `interest_amount_paid` DECIMAL(10, 2) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);
 
-- loan_delinquent_histories
CREATE TABLE `loan_delinquent_histories` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `loan_transaction_id` BIGINT NOT NULL COMMENT "refer to loan_transactions.id",
 `bills` JSON NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);
 
--  loan_payments
CREATE TABLE `loan_payments` (
 `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
 `loan_transaction_id` BIGINT NOT NULL COMMENT "refer to loan_transactions.id",
 `principal_amount` DECIMAL(10, 2) NOT NULL,
 `principal_amount_paid` DECIMAL(10, 2) NOT NULL,
 `interest_amount` DECIMAL(10, 2) NOT NULL,
 `interest_amount_paid` DECIMAL(10, 2) NOT NULL,
 `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 `created_by` VARCHAR(255) NOT NULL,
 `updated_at` TIMESTAMP NULL,
 `updated_by` VARCHAR(255) NULL,
 `deleted_at` TIMESTAMP NULL,
 `deleted_by` VARCHAR(255) NULL,
 `is_deleted` TINYINT NOT NULL DEFAULT 0
);

ALTER TABLE `users` ADD COLUMN `delinquent_level` INT NOT NULL DEFAULT 0 AFTER `password`;