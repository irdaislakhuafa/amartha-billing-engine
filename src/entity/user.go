package entity

type (
	CreateUserParams struct {
		Name     string `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required,max=255,min=1"`
		Email    string `db:"email" json:"email" form:"email" query:"email" params:"email" validate:"required,email"`
		Password string `db:"password" json:"password" form:"password" query:"password" params:"password" validate:"required,min=8,max=255"`
	}
	RegisterUserParams CreateUserParams

	GetUserParams struct {
		ID        int64  `json:"id" query:"id" form:"id" params:"id" validate:"required"`
		IsDeleted int    `json:"is_deleted" query:"is_deleted" form:"is_deleted" params:"is_deleted" validate:""`
		Email     string `json:"email" query:"email" form:"email" params:"email" validate:"omitempty,email"`
	}

	ListUserParams struct {
		PaginationParams
		IsDeleted int    `json:"is_deleted" query:"is_deleted" form:"is_deleted" params:"is_deleted" validate:""`
		Search    string `json:"search" query:"search" form:"search" params:"search" validate:""`
	}

	UpdateUserParams struct {
		Name  string `db:"name" json:"name" form:"name" query:"name" params:"name" validate:"required,max=255,min=1"`
		Email string `db:"email" json:"email" form:"email" query:"email" params:"email" validate:"required,email"`
		// Password string `db:"password" json:"password" form:"password" query:"password" params:"password" validate:"required,min=8,max=255"`
		ID int64 `db:"id" json:"id" form:"id" query:"id" params:"id" validate:"required"`
	}

	DeleteUserParams struct {
		ID        int64 `json:"id" form:"id" query:"id" params:"id" validate:"required"`
		IsDeleted int   `json:"is_deleted" form:"is_deleted" query:"is_deleted" params:"is_deleted" validate:""`
	}

	LoginUserParams struct {
		Email    string `db:"email" form:"email" query:"email" params:"email" json:"email" validate:"required,email"`
		Password string `db:"password" form:"password" query:"password" params:"password" json:"password" validate:"required,min=8,max=255"`
	}

	User struct {
		Name            string `db:"name" json:"name"`
		Email           string `db:"email" json:"email"`
		Password        string `db:"password" json:"password"`
		DelinquentLevel int    `db:"delinquent_level" json:"delinquent_level"`
		Base
	}
)
