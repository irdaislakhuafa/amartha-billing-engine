package usecase

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/usecase/loan"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/go-sdk/caches"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/smtp"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

type (
	Usecase struct {
		Loan loan.Interface
	}
)

func Init(
	log log.Interface,
	cfg config.Config,
	val *validator.Validate,
	db *sql.DB,
	dom *domain.Domain,
	smtp smtp.Interface,
	storage storage.Interface,
	cache caches.Interface[entity.Cache],
) *Usecase {
	return &Usecase{
		Loan: loan.Init(cfg, log, val, db, dom),
	}
}
