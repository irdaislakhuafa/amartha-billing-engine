package entity

import "github.com/shopspring/decimal"

type (
	CreateLoanParams struct {
		Name              string  `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required,min=1,max=255"`
		Description       string  `db:"description" json:"description" form:"description" query:"description" params:"description" validate:""`
		InterestRate      float64 `db:"interest_rate" json:"interest_rate" form:"interest_rate" query:"interest_rate" params:"interest_rate" validate:"gte=0"`
		RepaymentType     string  `db:"repayment_type" json:"repayment_type" form:"repayment_type" query:"repayment_type" params:"repayment_type" validate:"required,oneof=weeks months years"`
		RepaymentDuration int     `db:"repayment_duration" json:"repayment_duration" form:"repayment_duration" query:"repayment_duration" params:"repayment_duration" validate:"gte=1"`
	}

	ListLoanParams struct {
		IDs            []int64  `db:"ids" json:"ids" form:"ids" query:"ids" params:"ids" validate:"dive,required"`
		RepaymentTypes []string `db:"repayment_types" json:"repayment_types" form:"repayment_types" query:"repayment_types" params:"repayment_types" validate:"dive,oneof=weeks months years"`
		Search         string   `db:"search" json:"search" form:"search" query:"search" params:"search" validate:"max=255"`
		IsDeleted      int      `db:"is_deleted" json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted" validate:"oneof=0 1"`
		PaginationParams
	}

	GetLoanParams struct {
		ID        int64 `db:"id" json:"id" form:"id" query:"id" params:"id" validate:"required"`
		IsDeleted int   `db:"is_deleted" form:"is_deleted" query:"id" params:"id" validate:""`
	}

	UpdateLoanParams struct {
		Name              string  `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required,min=1,max=255"`
		Description       string  `db:"description" json:"description" form:"description" query:"description" params:"description" validate:""`
		InterestRate      float64 `db:"interest_rate" json:"interest_rate" form:"interest_rate" query:"interest_rate" params:"interest_rate" validate:"gte=0"`
		RepaymentType     string  `db:"repayment_type" json:"repayment_type" form:"repayment_type" query:"repayment_type" params:"repayment_type" validate:"required,oneof=weeks months years"`
		RepaymentDuration int     `db:"repayment_duration" json:"repayment_duration" form:"repayment_duration" query:"repayment_duration" params:"repayment_duration" validate:"gte=1"`
		ID                int64   `db:"id" json:"id" form:"id" query:"id" params:"id" validate:"required"`
	}

	DeleteLoanParams struct {
		ID        int64 `db:"id" json:"id" form:"id" query:"id" params:"id" validate:"required"`
		IsDeleted int   `db:"is_deleted" json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted" validate:"oneof=0 1"`
	}

	Loan struct {
		Name        string `db:"name" json:"name"`
		Description string `db:"description" json:"description"`
		// per annum
		InterestRate decimal.Decimal `db:"interest_rate" json:"interest_rate"`
		// weeks, months, years
		RepaymentType     string `db:"repayment_type" json:"repayment_type"`
		RepaymentDuration int32  `db:"repayment_duration" json:"repayment_duration"`
		Base
	}
)

const (
	REPAYMENT_TYPE_WEEKS  = "weeks"
	REPAYMENT_TYPE_MONTHS = "months"
	REPAYMENT_TYPE_YEARS  = "years"
)
