package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/errmessages"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/validation"
	"github.com/irdaislakhuafa/go-sdk/auth"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/cryptography"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
)

type (
	Interface interface {
		Register(ctx context.Context, params entity.RegisterUserParams) (entity.User, string, error)
		Get(ctx context.Context, params entity.GetUserParams) (entity.User, error)
		List(ctx context.Context, params entity.ListUserParams) ([]entity.User, entity.Pagination, error)
		Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error)
		Delete(ctx context.Context, params entity.DeleteUserParams) (entity.User, error)
		Login(ctx context.Context, params entity.LoginUserParams) (entity.User, string, error)
	}

	impl struct {
		log log.Interface
		db  *sql.DB
		dom *domain.Domain
		val *validator.Validate
		cfg config.Config
	}
)

func Init(
	log log.Interface,
	db *sql.DB,
	dom *domain.Domain,
	val *validator.Validate,
	cfg config.Config,
) Interface {
	return &impl{
		log: log,
		db:  db,
		dom: dom,
		val: val,
		cfg: cfg,
	}
}

func (i *impl) Register(ctx context.Context, params entity.RegisterUserParams) (entity.User, string, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.User{}, "", errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	_, err := i.dom.User.Get(ctx, entity.GetUserParams{
		Email:     params.Email,
		IsDeleted: 0,
	})
	code := errors.GetCode(err)
	if code.IsNotOneOf(codes.CodeSQLRecordDoesNotExist, codes.NoCode) {
		return entity.User{}, "", errors.NewWithCode(code, err.Error())
	}
	if code == codes.NoCode {
		return entity.User{}, "", errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_ALREADY_REGISTERED)
	}

	if pass, err := cryptography.NewBcrypt().Hash([]byte(params.Password)); err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	} else {
		params.Password = string(pass)
	}

	user, err := i.dom.User.Create(ctx, entity.CreateUserParams(params))
	if err != nil {
		return entity.User{}, "", errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	user.Password = ""
	return user, errmessages.USER_SUCCESS_REGISTERED, nil
}

func (i *impl) Get(ctx context.Context, params entity.GetUserParams) (entity.User, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	user, err := i.dom.User.Get(ctx, params)
	if err != nil {
		if errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
			return entity.User{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_NOT_REGISTERED)
		}
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	user.Password = ""
	return user, nil
}

func (i *impl) List(ctx context.Context, params entity.ListUserParams) ([]entity.User, entity.Pagination, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	results, p, err := i.dom.User.List(ctx, params)
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	for i := range results {
		results[i].Password = ""
	}

	return results, p, nil
}

func (i *impl) Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.User.Get(ctx, entity.GetUserParams{ID: params.ID, IsDeleted: 0})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.User{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_NOT_REGISTERED)
		}
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	up, err := i.dom.User.Update(ctx, params)
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	up.Password = ""
	up.CreatedAt = prev.CreatedAt
	up.CreatedBy = prev.CreatedBy
	return up, nil
}

func (i *impl) Delete(ctx context.Context, params entity.DeleteUserParams) (entity.User, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev, err := i.dom.User.Get(ctx, entity.GetUserParams{
		ID:        params.ID,
		IsDeleted: 0,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.User{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_NOT_REGISTERED)
		}
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	del, err := i.dom.User.Delete(ctx, params)
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	prev.Password = ""
	prev.DeletedAt = del.DeletedAt
	prev.DeletedBy = del.DeletedBy
	prev.IsDeleted = del.IsDeleted

	return prev, nil
}

func (i *impl) Login(ctx context.Context, params entity.LoginUserParams) (entity.User, string, error) {
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.User{}, "", errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	user, err := i.dom.User.Get(ctx, entity.GetUserParams{
		Email: params.Email,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.User{}, "", errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_NOT_REGISTERED)
		}
		return entity.User{}, "", errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	err = cryptography.NewBcrypt().Compare([]byte(params.Password), []byte(user.Password))
	if err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_PASSWORD_NOT_MATCH)
	}

	stoken, err := auth.
		InitJWT([]byte(i.cfg.Token.Secret), entity.JWTClaims{
			UID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(i.cfg.Token.ExpirationMinutes) * time.Minute)),
			},
		}).
		Generate(ctx)
	if err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeInternalServerError, err.Error())
	}

	user.Password = ""
	return user, stoken, nil
}
