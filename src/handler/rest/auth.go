package rest

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/ctxkey"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/errmessages"
	"github.com/irdaislakhuafa/go-sdk/auth"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/header"
)

func (r *rest) authJWT(c *fiber.Ctx) error {
	prefix := "Bearer "

	hAuth := c.Get(header.KeyAuthorization, "")
	hAuth = strings.Trim(hAuth, " ")
	isUnauth := hAuth == "" || len(hAuth) <= len(prefix)
	if isUnauth {
		return r.httpResError(c, errors.NewWithCode(codes.CodeUnauthorized, errmessages.AUTH_HEADER_REQUIRED))
	}

	hToken := hAuth[len(prefix):]
	jAuth := auth.InitJWT([]byte(r.cfg.Token.Secret), &entity.JWTClaims{})
	jToken, err := jAuth.Validate(c.UserContext(), hToken)
	if err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeUnauthorized, err.Error()))
	}

	jClaims, err := jAuth.ExtractClaims(c.UserContext(), jToken)
	if err != nil {
		return r.httpResError(c, errors.NewWithCode(codes.CodeUnauthorized, err.Error()))
	}

	ctx := c.UserContext()
	ctx = context.WithValue(ctx, ctxkey.USER_ID, jClaims.UID)
	ctx = context.WithValue(ctx, ctxkey.JWT_CLAIMS, jClaims)

	c.SetUserContext(ctx)

	return c.Next()
}
