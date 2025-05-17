package domain

import (
	"database/sql"

	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loan"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loanbilling"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loandelinquenthistories"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loanpayment"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/loantransaction"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain/user"
	entitygen "github.com/irdaislakhuafa/amartha-billing-engine/src/entity/gen"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/storage"
)

type (
	Domain struct {
		Loan                  loan.Interface
		User                  user.Interface
		LoanPayment           loanpayment.Interface
		LoanTransaction       loantransaction.Interface
		LoanBilling           loanbilling.Interface
		LoanDelinquentHistory loandelinquenthistories.Interface
	}
)

func Init(log log.Interface, queries *entitygen.Queries, db *sql.DB, storage storage.Interface) *Domain {
	return &Domain{
		Loan:                  loan.Init(queries, db, log),
		User:                  user.Init(log, queries),
		LoanPayment:           loanpayment.Init(log, queries),
		LoanTransaction:       loantransaction.Init(log, queries),
		LoanBilling:           loanbilling.Init(log, queries),
		LoanDelinquentHistory: loandelinquenthistories.Init(log, queries),
	}
}
