package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type (
	CreateLoanTransactionParams struct {
		InvoiceNumber string          `db:"invoice_number" json:"invoice_number" form:"invoice_number" params:"invoice_number" query:"invoice_number" validate:""`
		Notes         string          `db:"notes" json:"notes" form:"notes" params:"notes" query:"notes" validate:""`
		UserID        int64           `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
		User          json.RawMessage `db:"user" json:"user" form:"user" params:"user" query:"user" validate:""`
		LoanID        int64           `db:"loan_id" json:"loan_id" form:"loan_id" params:"loan_id" query:"loan_id" validate:"required"`
		Loan          json.RawMessage `db:"loan" json:"loan" form:"loan" params:"loan" query:"loan" validate:""`
		Amount        float64         `db:"amount" json:"amount" form:"amount" params:"amount" query:"amount" validate:"required,gte=0"`
	}

	UpdateLoanTransactionParams struct {
		InvoiceNumber string          `db:"invoice_number" json:"invoice_number" form:"invoice_number" params:"invoice_number" query:"invoice_number" validate:""`
		Notes         string          `db:"notes" json:"notes" form:"notes" params:"notes" query:"notes" validate:""`
		UserID        int64           `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
		User          json.RawMessage `db:"user" json:"user" form:"user" params:"user" query:"user" validate:""`
		LoanID        int64           `db:"loan_id" json:"loan_id" form:"loan_id" params:"loan_id" query:"loan_id" validate:"required"`
		Loan          json.RawMessage `db:"loan" json:"loan" form:"loan" params:"loan" query:"loan" validate:""`
		Amount        float64         `db:"amount" json:"amount" form:"amount" params:"amount" query:"amount" validate:"required"`
		ID            int64           `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	DeleteLoanTransactionParams struct {
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
	}

	ListLoanTransactionParams struct {
		Invoices     []string `db:"invoices" json:"invoices" form:"invoices" params:"invoices" query:"invoices"`
		UserIDs      []int64  `db:"user_ids" json:"user_ids" form:"user_ids" params:"user_ids" query:"user_ids"`
		LoanIDs      []int64  `db:"loan_ids" json:"loan_ids" form:"loan_ids" params:"loan_ids" query:"loan_ids"`
		IsDeleted    int8     `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted"`
		WithPayments bool     `db:"with_payments" json:"with_payments" form:"with_payments" params:"with_payments" query:"with_payments"`
		PaginationParams
	}

	GetLoanTransactionParams struct {
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		UserID    int64 `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:""`
	}

	CalculateOutstandingLoanTransactionParams struct {
		UserID int64 `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
	}

	CalculateOutstandingLoanTransaction struct {
		CurrentBillDate       *time.Time      `db:"current_bill_date" json:"current_bill_date"`
		NextBillDate          *time.Time      `db:"next_bill_date" json:"next_bill_date"`
		BilledPrincipalAmount decimal.Decimal `db:"billed_principal_amount" json:"billed_principal_amount"`
		BilledInterestAmount  decimal.Decimal `db:"billed_interest_amount" json:"billed_interest_amount"`
		TotalBilledAmount     decimal.Decimal `db:"total_billed_amount" json:"total_billed_amount"`
		TotalPaidAmount       decimal.Decimal `db:"total_paid_amount" json:"total_paid_amount"`
		OSPrincipalAmount     decimal.Decimal `db:"os_principal_amount" json:"os_principal_amount"`
		OSInterestAmount      decimal.Decimal `db:"os_interest_amount" json:"os_interest_amount"`
		TotalOSAmount         decimal.Decimal `db:"total_os_amount" json:"total_os_amount"`
		ListBilledBilling     []LoanBilling   `db:"list_billed_billing" json:"list_billed_billing"`
	}

	PayLoanTransactionParams struct {
		LoanTransactionID int64 `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		// PrincipalAmount   decimal.Decimal `db:"principal_amount" json:"principal_amount" form:"principal_amount" params:"principal_amount" query:"principal_amount" validate:"required"`
		// InterestAmount    decimal.Decimal `db:"interest_amount" json:"interest_amount" form:"interest_amount" params:"interest_amount" query:"interest_amount" validate:"required"`
		Amount float64 `db:"amount" json:"amount" form:"amount" params:"amount" query:"amount" validate:"required,gt=0"`
		UserID int64   `db:"user_id" json:"user_id" form:"user_id" params:"user_id" query:"user_id" validate:"required"`
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

		// related entity
		LoanBillings []LoanBilling `db:"loan_billings" json:"loan_billings,omitempty"`
		LoanPayments []LoanPayment `db:"loan_payment" json:"loan_payment,omitempty"`
	}
)

func (lt *LoanTransaction) GenInvoiceNumber(id int64, userID int64) string {
	lt.InvoiceNumber = fmt.Sprintf("LOAN/%d/%d", id, userID)
	return lt.InvoiceNumber
}

const (
	LOAN_TRANSACTION_STATUS_PAID   = "paid"
	LOAN_TRANSACTION_STATUS_UNPAID = "unpaid"
)
