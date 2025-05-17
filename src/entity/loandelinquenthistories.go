package entity

import (
	"encoding/json"
)

type (
	CreateLoanDelinquentHistoryParams struct {
		LoanTransactionID int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		Bills             json.RawMessage `db:"bills" json:"bills" form:"bills" params:"bills" query:"bills" validate:"required"`
	}

	UpdateLoanDelinquentHistoryParams struct {
		LoanTransactionID int64           `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:"required"`
		Bills             json.RawMessage `db:"bills" json:"bills" form:"bills" params:"bills" query:"bills" validate:"required"`
		ID                int64           `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	DeleteLoanDelinquentHistoryParams struct {
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
	}

	GetLoanDelinquentHistoryParams struct {
		ID        int64 `db:"id" json:"id" form:"id" params:"id" query:"id" validate:"required"`
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
	}

	ListLoanDelinquentHistoryParams struct {
		LoanTransactionID int64 `db:"loan_transaction_id" json:"loan_transaction_id" form:"loan_transaction_id" params:"loan_transaction_id" query:"loan_transaction_id" validate:""`
		IsDeleted         int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" params:"is_deleted" query:"is_deleted" validate:""`
		PaginationParams
	}

	LoanDelinquentHistory struct {
		// refer to loan_transactions.id
		LoanTransactionID int64           `db:"loan_transaction_id" json:"loan_transaction_id"`
		Bills             json.RawMessage `db:"bills" json:"bills"`
		Base
	}
)
