package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/go-sdk/appcontext"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/header"
	"github.com/irdaislakhuafa/go-sdk/language"
	"github.com/irdaislakhuafa/go-sdk/strformat"
)

func (r *rest) addFieldsToContext(c *fiber.Ctx) error {
	var ctx context.Context = c.Context()
	ctx = appcontext.SetAcceptLanguage(ctx, language.English)
	ctx = appcontext.SetRequestID(ctx, c.Get(header.KeyRequestID, uuid.New().String()))
	ctx = appcontext.SetUserAgent(ctx, c.Get(header.KeyUserAgent))
	ctx = appcontext.SetServiceVersion(ctx, r.cfg.Meta.Version)

	c.SetUserContext(ctx)
	return c.Next()
}

func (r *rest) setTimeout(c *fiber.Ctx) (err error) {
	// wrap context with timeout
	to := time.Duration(time.Second * time.Duration(r.cfg.Fiber.TimeoutSeconds))
	ctx, cancel := context.WithTimeout(c.UserContext(), to)
	defer cancel()
	defer func() {
		// if timeout is exceed then abort request and response error
		if ctx.Err() == context.DeadlineExceeded {
			err = r.httpResError(c, errors.NewWithCode(codes.CodeContextDeadlineExceeded, "Context deadline exceeded"))
			return
		}
	}()
	c.SetUserContext(ctx)
	err = c.Next()
	return err
}

func (r *rest) httpResSuccess(c *fiber.Ctx, code codes.Code, data any, p *entity.Pagination) error {
	ctx := c.UserContext()
	successApp := codes.Compile(code, appcontext.GetAcceptLanguage(ctx))
	res := entity.HTTPRes{
		Message: entity.HTTPMessage{
			Title: successApp.Title,
			Body:  successApp.Body,
		},
		Meta: entity.Meta{
			Path:       c.Path(),
			StatusCode: successApp.StatusCode,
			StatusStr:  http.StatusText(successApp.StatusCode),
			Message: strformat.TWE("{{ .Method }} {{ .URI }} {{ .StatusCode }} {{ .StatusStr }}", map[string]any{
				"Method":     c.Method(),
				"URI":        c.Context().URI(),
				"StatusCode": successApp.StatusCode,
				"StatusStr":  http.StatusText(successApp.StatusCode),
			}),
			Timestamp: time.Now().Format(time.RFC3339),
			Error:     nil,
			RequestID: appcontext.GetRequestID(ctx),
		},
		Data:       data,
		Pagination: p,
	}

	return c.Status(successApp.StatusCode).JSON(res, header.ContentTypeJSON)
}

func (r *rest) httpResError(c *fiber.Ctx, err error) error {
	ctx := c.UserContext()
	httpStatusCode, displayErr := errors.Compile(err, appcontext.GetAcceptLanguage(ctx))
	statusStr := http.StatusText(httpStatusCode)
	res := entity.HTTPRes{
		Message: entity.HTTPMessage{
			Title: displayErr.Title,
			Body:  displayErr.Body,
		},
		Meta: entity.Meta{
			Path:       c.Path(),
			StatusCode: httpStatusCode,
			StatusStr:  statusStr,
			Message: strformat.TWE("{{ .Method }} {{ .URI }} {{ .StatusCode }} {{ .StatusStr }}", map[string]any{
				"Method":     c.Method(),
				"URL":        c.Context().URI(),
				"StatusCode": httpStatusCode,
				"StatusStr":  statusStr,
			}),
			Timestamp: time.Now().Format(time.RFC3339),
			Error: &entity.MetaError{
				Code:    int(displayErr.Code),
				Message: err.Error(),
			},
			RequestID: appcontext.GetRequestID(ctx),
		},
		Data:       nil,
		Pagination: nil,
	}

	r.log.Error(ctx, err)
	c.Response().Header.Add(header.KeyRequestID, appcontext.GetRequestID(ctx))
	return c.Status(httpStatusCode).JSON(res, header.ContentTypeJSON)
}
