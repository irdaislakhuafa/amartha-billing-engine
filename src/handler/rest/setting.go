package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func (r *rest) CreateSetting(c *fiber.Ctx) error {
	body := entity.CreateSettingParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	setting, err := r.uc.Setting.Create(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, setting, nil)
}

func (r *rest) GetSetting(c *fiber.Ctx) error {
	params := entity.GetSettingParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	setting, err := r.uc.Setting.Get(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, setting, nil)
}

func (r *rest) ListSetting(c *fiber.Ctx) error {
	queries := entity.ListSettingParams{}
	if err := c.QueryParser(&queries); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	settings, pagination, err := r.uc.Setting.List(c.UserContext(), queries)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, settings, &pagination)
}

func (r *rest) UpdateSetting(c *fiber.Ctx) error {
	body := entity.UpdateSettingParams{}
	if err := c.BodyParser(&body); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	setting, err := r.uc.Setting.Update(c.UserContext(), body)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, setting, nil)
}

func (r *rest) DeleteSetting(c *fiber.Ctx) error {
	params := entity.DeleteSettingParams{}
	if err := c.ParamsParser(&params); err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeBadRequest, err.Error()))
	}

	setting, err := r.uc.Setting.Delete(c.UserContext(), params)
	if err != nil {
		return r.httpResError(c, err)
	}

	return r.httpResSuccess(c, codes.CodeSuccess, setting, nil)
}
