package domain

import (
	"database/sql"

	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/todo"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

type (
	Domain struct {
		Todo todo.Interface
	}
)

func Init(log log.Interface, queries *entitygen.Queries, db *sql.DB, storage storage.Interface) *Domain {
	return &Domain{
		Todo: todo.Init(log, queries, db),
	}
}
