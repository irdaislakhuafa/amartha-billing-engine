-- LOANS
-- name: CreateLoan :execresult
INSERT INTO `loans` (
  `name`, 
  `description`, 
  `interest_rate`, 
  `repayment_type`, 
  `repayment_duration`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetLoan :one
SELECT * FROM `loans`;

-- name: ListLoan :many
SELECT * FROM `loans`;

-- name: CountLoan :one
SELECT COUNT(`id`) AS `total` FROM `loans`;

-- name: UpdateLoan :execresult
UPDATE `loans` SET
  `name` = ?,
  `description` = ?,
  `interest_rate` = ?,
  `repayment_type` = ?,
  `repayment_duration` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteLoan :execresult
UPDATE `loans` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;

-- USERS
-- name: CreateUser :execresult
INSERT INTO `users` (
  `name`, 
  `email`, 
  `password`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT * FROM `users`;

-- name: ListUser :many
SELECT * FROM `users`;

-- name: CountUser :one
SELECT COUNT(`id`) AS `total` FROM `users`;

-- name: UpdateUser :execresult
UPDATE `users` SET
  `name` = ?,
  `email` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteUser :execresult
UPDATE `users` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;

-- LOAN TRANSACTIONS
-- name: CreateLoanTransaction :execresult
INSERT INTO `loan_transactions` (
  `invoice_number`, 
  `notes`, 
  `user_id`, 
  `user`, 
  `loan_id`, 
  `loan`, 
  `amount`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetLoanTransaction :one
SELECT * FROM `loan_transactions`;

-- name: ListLoanTransaction :many
SELECT * FROM `loan_transactions`;

-- name: CountLoanTransaction :one
SELECT COUNT(`id`) AS `total` FROM `loan_transactions`;

-- name: UpdateLoanTransaction :execresult
UPDATE `loan_transactions` SET
  `invoice_number` = ?,
  `notes` = ?,
  `user_id` = ?,
  `user` = ?,
  `loan_id` = ?,
  `loan` = ?,
  `amount` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteLoanTransaction :execresult
UPDATE `loan_transactions` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;

-- LOAN BILLING
-- name: CreateLoanBilling :execresult
INSERT INTO `loan_billings` (
  `loan_transaction_id`, 
  `bill_date`, 
  `principal_amount`, 
  `principal_amount_paid`, 
  `interest_amount`, 
  `interest_amount_paid`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetLoanBilling :one
SELECT * FROM `loan_billings`;

-- name: ListLoanBilling :many
SELECT * FROM `loan_billings`;

-- name: UpdateLoanBilling :execresult
UPDATE `loan_billings` SET
  `loan_transaction_id` = ?,
  `bill_date` = ?,
  `principal_amount` = ?,
  `principal_amount_paid` = ?,
  `interest_amount` = ?,
  `interest_amount_paid` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteLoanBilling :execresult
UPDATE `loan_billings` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;

-- LOAN DELINQUENT HISTORY
-- name: CreateLoanDelinquentHistory :execresult
INSERT INTO `loan_delinquent_histories` (
  `loan_transaction_id`, 
  `bills`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?);

-- name: GetLoanDelinquentHistory :one
SELECT * FROM `loan_delinquent_histories`;

-- name: ListLoanDelinquentHistory :many
SELECT * FROM `loan_delinquent_histories`;

-- name: UpdateLoanDelinquentHistory :execresult
UPDATE `loan_delinquent_histories` SET
  `loan_transaction_id` = ?,
  `bills` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteLoanDelinquentHistory :execresult
UPDATE `loan_delinquent_histories` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;

-- LOAN PAYMENTS
-- name: CreateLoanPayment :execresult
INSERT INTO `loan_payments` (
  `loan_transaction_id`, 
  `principal_amount`, 
  `principal_amount_paid`, 
  `interest_amount`, 
  `interest_amount_paid`, 
  `created_at`, 
  `created_by`
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetLoanPayment :one
SELECT * FROM `loan_payments`;

-- name: ListLoanPayment :many
SELECT * FROM `loan_payments`;

-- name: UpdateLoanPayment :execresult
UPDATE `loan_payments` SET
  `loan_transaction_id` = ?,
  `principal_amount` = ?,
  `principal_amount_paid` = ?,
  `interest_amount` = ?,
  `interest_amount_paid` = ?,
  `updated_at` = ?,
  `updated_by` = ?
WHERE `id` = ?;

-- name: DeleteLoanPayment :execresult
UPDATE `loan_payments` SET
  `is_deleted` = ?,
  `deleted_at` = ?,
  `deleted_by` = ?
WHERE `id` = ?;
