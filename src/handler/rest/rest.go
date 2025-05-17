package rest

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/usecase"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/operator"
)

var once = &sync.Once{}

type (
	Interface interface {
		Run()
	}
	rest struct {
		cfg config.Config
		svr *fiber.App
		log log.Interface
		uc  *usecase.Usecase
	}
)

func Init(cfg config.Config, log log.Interface, uc *usecase.Usecase) Interface {
	r := &rest{}
	once.Do(func() {
		preforks := map[string]bool{
			config.APP_MODE_DEV:  false,
			config.APP_MODE_PROD: true,
		}
		svr := fiber.New(fiber.Config{
			Prefork:               preforks[cfg.Fiber.Mode], // enable multi processes for better performance
			DisableStartupMessage: true,
		})

		r = &rest{
			cfg: cfg,
			svr: svr,
			log: log,
			uc:  uc,
		}

		// if prefer default then ignore custom cors
		if cfg.Fiber.Cors.PreferDefault {
			svr.Use(cors.New())
		} else { // if PreferDefault is disabled, then use custom cors
			svr.Use(cors.New(cors.Config{
				AllowOrigins:     strings.Join(cfg.Fiber.Cors.Origins, ", "),
				AllowMethods:     strings.Join(cfg.Fiber.Cors.Methods, ", "),
				AllowHeaders:     strings.Join(cfg.Fiber.Cors.Headers, ", "),
				AllowCredentials: cfg.Fiber.Cors.Cookies,
			}))
		}

		// enable recovery on panic
		r.svr.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))

		// add fields to context
		r.svr.Use(r.addFieldsToContext)

		// set timeout
		r.svr.Use(r.setTimeout)

		// register router
		r.RegisterRoutes()
	})

	return r
}

func (r *rest) Run() {
	port := operator.Ternary(r.cfg.Fiber.Port != "", r.cfg.Fiber.Port, "8000")
	r.log.Info(context.Background(), fmt.Sprintf("Rest API listening at port '%s'...", port))
	if err := r.svr.Listen(":" + port); err != nil {
		r.log.Fatal(context.Background(), err)
	}
}

/**
 * @Summary Heatlh Check
 * @Description This endpoint will hit the server
 * @Tags server
 * @Produce JSON
 * @Success 200 string example="PONG!"
 * @Router /ping [GET]
 */
func (r *rest) Ping(c *fiber.Ctx) error {
	return r.httpResSuccess(c, codes.CodeSuccess, "PONG!", nil)
}
