package setting

import (
	"context"
	"database/sql"
	"time"

	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/ctxkey"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/querybuilder/sqlc"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateSettingParams) (entity.Setting, error)
		Get(ctx context.Context, params entity.GetSettingParams) (entity.Setting, error)
		Update(ctx context.Context, params entity.UpdateSettingParams) (entity.Setting, error)
		Delete(ctx context.Context, params entity.DeleteSettingParams) (entity.Setting, error)
		List(ctx context.Context, params entity.ListSettingParams) ([]entity.Setting, entity.Pagination, error)
	}

	impl struct {
		queries *entitygen.Queries
		log     log.Interface
	}
)

func Init(queries *entitygen.Queries, log log.Interface) Interface {
	return &impl{
		queries: queries,
		log:     log,
	}
}

// Create implements Interface.
func (i *impl) Create(ctx context.Context, params entity.CreateSettingParams) (entity.Setting, error) {
	args := entitygen.CreateSettingParams{
		Name:      params.Name,
		Value:     params.Value,
		CreatedAt: time.Now(),
		CreatedBy: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	setting, err := i.queries.CreateSetting(ctx, args)
	if err != nil {
		return entity.Setting{}, err
	}

	result := entity.Setting{
		Name:  params.Name,
		Value: params.Value,
		Base: entity.Base{
			CreatedAt: args.CreatedAt,
			CreatedBy: args.CreatedBy,
		},
	}

	if result.ID, err = setting.LastInsertId(); err != nil {
		return entity.Setting{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	return result, nil
}

// Delete implements Interface.
func (i *impl) Delete(ctx context.Context, params entity.DeleteSettingParams) (entity.Setting, error) {
	args := entitygen.DeleteSettingParams{
		IsDeleted: params.IsDeleted,
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		DeletedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:        params.ID,
	}

	_, err := i.queries.DeleteSetting(ctx, args)
	if err != nil {
		return entity.Setting{}, err
	}

	result := entity.Setting{
		Base: entity.Base{
			ID:        params.ID,
			IsDeleted: params.IsDeleted,
			DeletedAt: &args.DeletedAt.Time,
			DeletedBy: &args.DeletedBy.String,
		},
	}

	return result, nil
}

// Get implements Interface.
func (i *impl) Get(ctx context.Context, params entity.GetSettingParams) (entity.Setting, error) {
	setting, err := i.queries.GetSetting(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("(id = ? OR name = ?)", params.ID, params.Name)
		b.And("is_deleted = ?", params.IsDeleted)
	}))
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result, err := i.rowToEntity(setting)
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params entity.ListSettingParams) ([]entity.Setting, entity.Pagination, error) {
	paramsBackup := params
	if err := params.Parse(); err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	ctx = sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.And("is_deleted = ?", params.IsDeleted)
		if params.Search != "" {
			params.Search = "%" + params.Search + "%"
			b.And("(name LIKE ? OR value LIKE ?)", params.Search, params.Search)
		}
	})

	rows, err := i.queries.ListSetting(sqlc.Build(ctx, func(b *sqlc.Builder) {
		b.Limit(params.Limit)
		b.Offset(params.Page)
		b.Order(params.OrderBy + " " + params.OrderType)
	}))
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	var results []entity.Setting
	for _, row := range rows {
		result, err := i.rowToEntity(row)
		if err != nil {
			return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		results = append(results, result)
	}

	total, err := i.queries.CountSetting(ctx)
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	p := entity.GenPagination(paramsBackup.Page, paramsBackup.Limit, int(total), []string{params.OrderBy, params.OrderType})

	return results, p, nil
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateSettingParams) (entity.Setting, error) {
	args := entitygen.UpdateSettingParams{
		Name:      params.Name,
		Value:     params.Value,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy: sql.NullString{String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)), Valid: true},
		ID:        params.ID,
	}

	_, err := i.queries.UpdateSetting(ctx, args)
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	result := entity.Setting{
		Name:  params.Name,
		Value: params.Value,
		Base: entity.Base{
			ID:        params.ID,
			UpdatedAt: convert.SQLNullTimeToNil(args.UpdatedAt),
			UpdatedBy: convert.SQLNullStringToNil(args.UpdatedBy),
		},
	}

	return result, nil
}

func (i *impl) rowToEntity(row entitygen.Setting) (entity.Setting, error) {
	result := entity.Setting{
		Name:  row.Name,
		Value: row.Value,
		Base: entity.Base{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			CreatedBy: row.CreatedBy,
			UpdatedAt: convert.SQLNullTimeToNil(row.UpdatedAt),
			UpdatedBy: convert.SQLNullStringToNil(row.UpdatedBy),
			DeletedAt: convert.SQLNullTimeToNil(row.DeletedAt),
			DeletedBy: convert.SQLNullStringToNil(row.DeletedBy),
			IsDeleted: row.IsDeleted,
		},
	}

	return result, nil
}
