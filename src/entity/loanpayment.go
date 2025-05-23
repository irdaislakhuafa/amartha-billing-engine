package entity

import (
	"github.com/shopspring/decimal"
)

type (
	CreateLoanPaymentParams struct {
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount" form:"principal_amount" params:"principal_amount" query:"principal_amount" validate:"required"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid" form:"principal_amount_paid" params:"principal_amount_paid" query:"principal_amount_paid" validate:"required"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount" form:"interest_amount" params:"interest_amount" query:"interest_amount" validate:"required"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid" form:"interest_amount_paid" params:"interest_amount_paid" query:"interest_amount_paid" validate:"required"`
		LoanBillingID       int64           `db:"loan_billing_id" json:"loan_billing_id" form:"loan_billing_id" params:"loan_billing_id" query:"loan_billing_id" validate:"required"`
	}

	UpdateLoanPaymentParams struct {
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount" form:"principal_amount" params:"principal_amount" query:"principal_amount" validate:"required"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid" form:"principal_amount_paid" params:"principal_amount_paid" query:"principal_amount_paid" validate:"required"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount" form:"interest_amount" params:"interest_amount" query:"interest_amount" validate:"required"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid" form:"interest_amount_paid" params:"interest_amount_paid" query:"interest_amount_paid" validate:"required"`
		ID                  int64           `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	DeleteLoanPaymentParams struct {
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	ListLoanPaymentParams struct {
		IDs                    []int64          `db:"id" json:"id" form:"id" params:"id" query:"id" validate:""`
		LoanTransactionIDs     []int64          `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:""`
		PrincipalAmountGTE     *decimal.Decimal `db:"principal_amount_gte" json:"principal_amount_gte" form:"principal_amount_gte" params:"principal_amount_gte" query:"principal_amount_gte" validate:""`
		PrincipalAmountLTE     *decimal.Decimal `db:"principal_amount_lte" json:"principal_amount_lte" form:"principal_amount_lte" params:"principal_amount_lte" query:"principal_amount_lte" validate:""`
		PrincipalAmountPaidGTE *decimal.Decimal `db:"principal_amount_paid_gte" json:"principal_amount_paid_gte" form:"principal_amount_paid_gte" params:"principal_amount_paid_gte" query:"principal_amount_paid_gte" validate:""`
		PrincipalAmountPaidLTE *decimal.Decimal `db:"principal_amount_paid_lte" json:"principal_amount_paid_lte" form:"principal_amount_paid_lte" params:"principal_amount_paid_lte" query:"principal_amount_paid_lte" validate:""`
		InterestAmountGTE      *decimal.Decimal `db:"interest_amount_gte" json:"interest_amount_gte" form:"interest_amount_gte" params:"interest_amount_gte" query:"interest_amount_gte" validate:""`
		InterestAmountLTE      *decimal.Decimal `db:"interest_amount_lte" json:"interest_amount_lte" form:"interest_amount_lte" params:"interest_amount_lte" query:"interest_amount_lte" validate:""`
		InterestAmountPaidGTE  *decimal.Decimal `db:"interest_amount_paid_gte" json:"interest_amount_paid_gte" form:"interest_amount_paid_gte" params:"interest_amount_paid_gte" query:"interest_amount_paid_gte" validate:""`
		InterestAmountPaidLTE  *decimal.Decimal `db:"interest_amount_paid_lte" json:"interest_amount_paid_lte" form:"interest_amount_paid_lte" params:"interest_amount_paid_lte" query:"interest_amount_paid_lte" validate:""`
		IsDeleted              int8             `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		PaginationParams
	}

	LoanPayment struct {
		// refer to loan_transactions.id
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid"`
		LoanBillingID       int64           `db:"loan_billing_id" json:"loan_billing_id"`
		Base
	}
)
