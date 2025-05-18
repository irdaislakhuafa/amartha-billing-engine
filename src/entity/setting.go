package entity

type (
	CreateSettingParams struct {
		Name  string `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required"`
		Value string `db:"value" json:"value" form:"value" query:"value" params:"value" validate:"required"`
	}

	UpdateSettingParams struct {
		Name  string `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required"`
		Value string `db:"value" json:"value" form:"value" query:"value" params:"value" validate:"required"`
		ID    int64  `db:"id" json:"id" form:"id" query:"id" params:"id" validate:"required"`
	}

	DeleteSettingParams struct {
		IsDeleted int8  `db:"is_deleted" json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted"`
		ID        int64 `db:"id" json:"id" form:"id" query:"id" params:"id"`
	}

	GetSettingParams struct {
		ID        int64  `db:"id" json:"id" form:"id" query:"id" params:"id"`
		Name      string `db:"name" json:"name" form:"name" query:"name" params:"name"`
		IsDeleted int8   `db:"is_deleted" json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted"`
	}

	ListSettingParams struct {
		PaginationParams
		Search    string `db:"search" json:"search" form:"search" query:"search" params:"search"`
		IsDeleted int8   `db:"is_deleted" json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted"`
	}

	Setting struct {
		Name  string `db:"name" json:"name" form:"name" query:"name" params:"name"`
		Value string `db:"value" json:"value" form:"value" query:"value" params:"value"`
		Base
	}
)

const (
	SETTING_NAME_EOD_DATE                     = "eod_date"
	SETTING_NAME_LIMIT_BILLING_FOR_DELINQUENT = "limit_billing_for_delinquent"
)
