package entity

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type (
	CreateLoanTransactionParams struct {
		InvoiceNumber string          `db:"invoice_number" json:"invoice_number" form:"invoice_number" params:"invoice_number" query:"invoice_number" validate:""`
		Notes         string          `db:"notes" json:"notes" form:"notes" params:"notes" query:"notes" validate:"required"`
		UserID        int64           `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
		User          json.RawMessage `db:"user" json:"user" form:"user" params:"user" query:"user" validate:""`
		LoanID        int64           `db:"loan_id" json:"loan_id" form:"loan_id" params:"loan_id" query:"loan_id" validate:"required"`
		Loan          json.RawMessage `db:"loan" json:"loan" form:"loan" params:"loan" query:"loan" validate:""`
		Amount        decimal.Decimal `db:"amount" json:"amount" form:"amount" params:"amount" query:"amount" validate:"required"`
	}

	UpdateLoanTransactionParams struct {
		InvoiceNumber string          `db:"invoice_number" json:"invoice_number" form:"invoice_number" params:"invoice_number" query:"invoice_number" validate:""`
		Notes         string          `db:"notes" json:"notes" form:"notes" params:"notes" query:"notes" validate:""`
		UserID        int64           `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
		User          json.RawMessage `db:"user" json:"user" form:"user" params:"user" query:"user" validate:""`
		LoanID        int64           `db:"loan_id" json:"loan_id" form:"loan_id" params:"loan_id" query:"loan_id" validate:"required"`
		Loan          json.RawMessage `db:"loan" json:"loan" form:"loan" params:"loan" query:"loan" validate:""`
		Amount        decimal.Decimal `db:"amount" json:"amount" form:"amount" params:"amount" query:"amount" validate:"required"`
		ID            int64           `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	DeleteLoanTransactionParams struct {
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
	}

	ListLoanTransactionParams struct {
		Invoices  []string `db:"invoices" json:"invoices" form:"invoices" params:"invoices" query:"invoices"`
		UserIDs   []int64  `db:"user_ids" json:"user_ids" form:"user_ids" params:"user_ids" query:"user_ids"`
		LoanIDs   []int64  `db:"loan_ids" json:"loan_ids" form:"loan_ids" params:"loan_ids" query:"loan_ids"`
		IsDeleted int8     `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted"`
		PaginationParams
	}

	GetLoanTransactionParams struct {
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
	}

	LoanTransaction struct {
		InvoiceNumber string `db:"invoice_number" json:"invoice_number"`
		Notes         string `db:"notes" json:"notes"`
		// refer to users.id
		UserID int64           `db:"user_id" json:"user_id"`
		User   json.RawMessage `db:"user" json:"user"`
		// refer to loans.id
		LoanID int64           `db:"loan_id" json:"loan_id"`
		Loan   json.RawMessage `db:"loan" json:"loan"`
		Amount decimal.Decimal `db:"amount" json:"amount"`
		Base
	}
)

func (lt *LoanTransaction) GenInvoiceNumber(id int64, userID int64) string {
	lt.InvoiceNumber = fmt.Sprintf("LOAN/%d/%d", id, userID)
	return lt.InvoiceNumber
}
