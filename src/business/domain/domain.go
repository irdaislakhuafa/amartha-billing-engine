package domain

import (
	"database/sql"

	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loan"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

type (
	Domain struct {
		Loan loan.Interface
	}
)

func Init(log log.Interface, queries *entitygen.Queries, db *sql.DB, storage storage.Interface) *Domain {
	return &Domain{
		Loan: loan.Init(queries, db, log),
	}
}
