package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	CreateLoanBillingParams struct {
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		BillDate            time.Time       `db:"bill_date" json:"bill_date" form:"bill_date" params:"bill_date" query:"bill_date" validate:"required"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount" form:"principal_amount" params:"principal_amount" query:"principal_amount" validate:"required"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid" form:"principal_amount_paid" params:"principal_amount_paid" query:"principal_amount_paid" validate:"required"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount" form:"interest_amount" params:"interest_amount" query:"interest_amount" validate:"required"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid" form:"interest_amount_paid" params:"interest_amount_paid" query:"interest_amount_paid" validate:"required"`
	}

	UpdateLoanBillingParams struct {
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		BillDate            time.Time       `db:"bill_date" json:"bill_date" form:"bill_date" params:"bill_date" query:"bill_date" validate:"required"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount" form:"principal_amount" params:"principal_amount" query:"principal_amount" validate:"required"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid" form:"principal_amount_paid" params:"principal_amount_paid" query:"principal_amount_paid" validate:"required"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount" form:"interest_amount" params:"interest_amount" query:"interest_amount" validate:"required"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid" form:"interest_amount_paid" params:"interest_amount_paid" query:"interest_amount_paid" validate:"required"`
		ID                  int64           `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	DeleteLoanBillingParams struct {
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	LoanBilling struct {
		// refer to loan_transactions.id
		LoanTransactionID   int64           `db:"loan_transaction_id" json:"loan_transaction_id"`
		BillDate            time.Time       `db:"bill_date" json:"bill_date"`
		PrincipalAmount     decimal.Decimal `db:"principal_amount" json:"principal_amount"`
		PrincipalAmountPaid decimal.Decimal `db:"principal_amount_paid" json:"principal_amount_paid"`
		InterestAmount      decimal.Decimal `db:"interest_amount" json:"interest_amount"`
		InterestAmountPaid  decimal.Decimal `db:"interest_amount_paid" json:"interest_amount_paid"`
		Base
	}
)
