package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/usecase"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/handler/rest"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/handler/scheduller"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/go-sdk/caches"
	"github.com/irdaislakhuafa/go-sdk/db"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/querybuilder/sqlc"
	"github.com/irdaislakhuafa/go-sdk/smtp"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

const (
	configFileJSON = "etc/cfg/conf.json"
)

func main() {
	ctx := context.Background()

	// read config file
	cfg, err := config.ReadFileJSON(configFileJSON)
	if err != nil {
		panic(err)
	}

	// initialize log
	l := log.Init(cfg.Log)

	// initialize db
	db, err := db.Init(cfg.DB.Master)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	l.Info(ctx, "Initialize db...")

	// initialize validator
	v := validator.New(
		validator.WithRequiredStructEnabled(),
	)
	l.Info(ctx, "Initialize validator...")

	// initialize storage client
	var s storage.Interface
	if cfg.Storage.AccessKeyID != "" && cfg.Storage.AccessKeySecret != "" {
		s, err = storage.Init(cfg.Storage)
		if err != nil {
			panic(err)
		}
	}
	l.Info(ctx, "Initialize storage...")

	// initialize cache
	c, err := caches.Init[entity.Cache](cfg.Cache)
	if err != nil {
		panic(err)
	}
	l.Info(ctx, "Initialize cache...")

	// initialize queries
	q := entitygen.New(sqlc.Wrap(db, sqlc.WrappedOpts{
		ShowQuery: true,
		ShowArgs:  true,
	}))
	l.Info(ctx, "Initialize query...")

	// initialize smtp
	smtp := smtp.InitGoMail(cfg.SMTP)
	l.Info(ctx, "Initialize smtp...")

	// initialize domain
	dom := domain.Init(l, q, db, s)
	l.Info(ctx, "Initialize domain...")

	// initialize usecase
	uc := usecase.Init(l, cfg, v, db, dom, smtp, s, c)
	l.Info(ctx, "Initialize usecase...")

	// initialize scheduller
	sch := scheduller.Init(l, cfg, uc)
	l.Info(ctx, "Initialize scheduller...")
	sch.Run()

	// choose running mode
	mode := flag.String("mode", "rest", "Please select running mode: [rest, gql, grpc]")
	flag.Parse()
	*mode = strings.ToLower(*mode)

	switch *mode {
	case "rest":
		r := rest.Init(cfg, l, uc)
		r.Run()
		// TODO: write for grpc and gql
	default:
		l.Fatal(context.Background(), fmt.Sprintf("Running mode '%s' not supported!", *mode))
	}
}
