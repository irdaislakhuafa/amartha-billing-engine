package setting

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/errmessages"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/validation"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
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
		dom *domain.Domain
		cfg config.Config
		log log.Interface
		db  *sql.DB
		val *validator.Validate
	}
)

func Init(log log.Interface, val *validator.Validate, cfg config.Config, db *sql.DB, dom *domain.Domain) Interface {
	return &impl{
		log: log,
		val: val,
		cfg: cfg,
		db:  db,
		dom: dom,
	}
}

// Create implements Interface.
func (i *impl) Create(ctx context.Context, params entity.CreateSettingParams) (entity.Setting, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Setting.Create(ctx, params)
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

// Delete implements Interface.
func (i *impl) Delete(ctx context.Context, params entity.DeleteSettingParams) (entity.Setting, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.Setting.Get(ctx, entity.GetSettingParams{
		ID:        params.ID,
		IsDeleted: 0,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.Setting{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_NOT_FOUND)
		}
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	del, err := i.dom.Setting.Delete(ctx, params)
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev.DeletedAt = del.DeletedAt
	prev.DeletedBy = del.DeletedBy
	prev.IsDeleted = del.IsDeleted

	return prev, nil
}

// Get implements Interface.
func (i *impl) Get(ctx context.Context, params entity.GetSettingParams) (entity.Setting, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Setting.Get(ctx, params)
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.Setting{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_NOT_FOUND)
		}
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, nil
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params entity.ListSettingParams) ([]entity.Setting, entity.Pagination, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return []entity.Setting{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, pagination, err := i.dom.Setting.List(ctx, params)
	if err != nil {
		return []entity.Setting{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return result, pagination, nil
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateSettingParams) (entity.Setting, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.Setting.Get(ctx, entity.GetSettingParams{
		ID:        params.ID,
		IsDeleted: 0,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.Setting{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_NOT_FOUND)
		}
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result, err := i.dom.Setting.Update(ctx, params)
	if err != nil {
		return entity.Setting{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result.CreatedAt = prev.CreatedAt
	result.CreatedBy = prev.CreatedBy

	return result, nil
}
