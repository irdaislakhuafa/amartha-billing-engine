package loantransaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/business/domain"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/entity"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/config"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/errmessages"
	"github.com/irdaislakhuafa/amartha-billing-engine/src/utils/validation"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/shopspring/decimal"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateLoanTransactionParams) (entity.LoanTransaction, error)
		List(ctx context.Context, params entity.ListLoanTransactionParams) ([]entity.LoanTransaction, entity.Pagination, error)
		Get(ctx context.Context, params entity.GetLoanTransactionParams) (entity.LoanTransaction, error)
		Update(ctx context.Context, params entity.UpdateLoanTransactionParams) (entity.LoanTransaction, error)
		Delete(ctx context.Context, params entity.DeleteLoanTransactionParams) (entity.LoanTransaction, error)
		CalculateOutstanding(ctx context.Context, params entity.CalculateOutstandingLoanTransactionParams) (entity.CalculateOutstandingLoanTransaction, error)
		WithTx(ctx context.Context, tx *sql.Tx) Interface
	}

	impl struct {
		log log.Interface
		val *validator.Validate
		cfg config.Config
		db  *sql.DB
		dom *domain.Domain
	}
)

func Init(log log.Interface, val *validator.Validate, cfg config.Config, db *sql.DB, dom *domain.Domain) Interface {
	return &impl{
		log: log,
		val: val,
		cfg: cfg,
		db:  db,
		dom: dom,
	}
}

// Create implements Interface.
func (i *impl) Create(ctx context.Context, params entity.CreateLoanTransactionParams) (entity.LoanTransaction, error) {
	// validate params
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	// get user
	var err error
	var user entity.User
	{
		user, err = i.dom.User.Get(ctx, entity.GetUserParams{
			ID:        params.UserID,
			IsDeleted: 0,
		})
		if err != nil {
			code := errors.GetCode(err)
			if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
				return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_NOT_REGISTERED)
			}
			return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		user.Password = ""
		userBytes, err := json.Marshal(user)
		if err != nil {
			return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeJSONMarshalError, err.Error())
		}
		params.User = userBytes
	}

	// get loan
	var loan entity.Loan
	{
		loan, err = i.dom.Loan.Get(ctx, entity.GetLoanParams{
			ID:        params.LoanID,
			IsDeleted: 0,
		})
		if err != nil {
			code := errors.GetCode(err)
			if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
				return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_NOT_FOUND)
			}
			return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}

		loanBytes, err := json.Marshal(loan)
		if err != nil {
			return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeJSONMarshalError, err.Error())
		}
		params.Loan = loanBytes
	}

	// ensure user has not been delinquent
	if user.DelinquentLevel > 0 {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.USER_IS_DELINQUENT)
	}

	// begin db tx
	tx, err := i.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	dLoanTransaction := i.dom.LoanTransaction.WithTx(ctx, tx)
	dLoanBilling := i.dom.LoanBilling.WithTx(ctx, tx)

	// create loan transaction
	result, err := dLoanTransaction.Create(ctx, params)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	// generate loan billing
	listBilling := []entity.LoanBilling{}
	switch loan.RepaymentType {
	case entity.REPAYMENT_TYPE_WEEKS:
		billDate := time.Now()
		for i := 0; i < int(loan.RepaymentDuration); i++ {
			listBilling = append(listBilling, entity.LoanBilling{
				LoanTransactionID: result.ID,
				BillDate:          billDate.AddDate(0, 0, 7),
			})
			billDate = billDate.AddDate(0, 0, 7)
		}
	case entity.REPAYMENT_TYPE_MONTHS:
		fallthrough
	case entity.REPAYMENT_TYPE_YEARS:
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_REPAYMENT_TYPE_NOT_AVAILABLE)
	default:
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_REPAYMENT_TYPE_INVALID)
	}

	// create loan billing
	basePrincipal := decimal.NewFromFloat(params.Amount).Div(decimal.NewFromInt(int64(loan.RepaymentDuration)))
	remainder := decimal.NewFromFloat(params.Amount).Mod(decimal.NewFromInt(int64(loan.RepaymentDuration)))

	createdBillings := []entity.LoanBilling{}
	for i, billing := range listBilling {
		principal := basePrincipal
		if i < int(remainder.IntPart()) {
			principal = principal.Add(remainder)
		}

		interest := principal.Div(decimal.NewFromInt(100)).Mul(loan.InterestRate)
		createdBilling, err := dLoanBilling.Create(ctx, entity.CreateLoanBillingParams{
			LoanTransactionID:   result.ID,
			UserID:              params.UserID,
			BillDate:            billing.BillDate,
			PrincipalAmount:     principal,
			PrincipalAmountPaid: decimal.NewFromInt(0),
			InterestAmount:      interest,
			InterestAmountPaid:  decimal.NewFromInt(0),
		})
		if err != nil {
			return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}

		createdBillings = append(createdBillings, createdBilling)
	}

	// commit tx
	if err := tx.Commit(); err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	result.LoanBilling = createdBillings
	return result, nil
}

func (i *impl) CalculateOutstanding(ctx context.Context, params entity.CalculateOutstandingLoanTransactionParams) (entity.CalculateOutstandingLoanTransaction, error) {
	// validate params
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.CalculateOutstandingLoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	result := entity.CalculateOutstandingLoanTransaction{
		CurrentBillDate:       nil,
		NextBillDate:          nil,
		BilledPrincipalAmount: decimal.Decimal{},
		BilledInterestAmount:  decimal.Decimal{},
		TotalBilledAmount:     decimal.Decimal{},
		TotalPaidAmount:       decimal.Decimal{},
		OSPrincipalAmount:     decimal.Decimal{},
		OSInterestAmount:      decimal.Decimal{},
		TotalOSAmount:         decimal.Decimal{},
	}

	// get billing
	billings, _, err := i.dom.LoanBilling.List(ctx, entity.ListLoanBillingParams{
		PaginationParams: entity.PaginationParams{
			Limit:     9_999_999,
			Page:      0,
			OrderBy:   "bill_date",
			OrderType: "desc",
		},
		IsDeleted: 0,
	})
	if err != nil {
		return result, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	if len(billings) > 0 {
		now := time.Now()
		for _, b := range billings {
			// get current billing
			isCurrent := b.BillDate.Before(now) && (b.PrincipalAmountPaid.LessThan(b.PrincipalAmount) || b.InterestAmountPaid.LessThan(b.InterestAmount))
			if isCurrent {
				result.CurrentBillDate = &b.BillDate
			}

			// get next billing
			isNext := b.BillDate.After(now) && (b.PrincipalAmountPaid.LessThan(b.PrincipalAmount) || b.InterestAmountPaid.LessThan(b.InterestAmount))
			if isNext {
				result.NextBillDate = &b.BillDate
			}

			// get billed principal amount
			if b.BillDate.Before(now) {
				result.BilledPrincipalAmount = result.BilledPrincipalAmount.Add(b.PrincipalAmount)
				result.BilledInterestAmount = result.BilledInterestAmount.Add(b.InterestAmount)
				result.TotalBilledAmount = result.TotalBilledAmount.Add(b.PrincipalAmount).Add(b.InterestAmount)
				result.TotalPaidAmount = result.TotalPaidAmount.Add(b.PrincipalAmountPaid).Add(b.InterestAmountPaid)
			}

			// get outstanding principal amount
			// if b.BillDate.Before(now) {
			result.OSPrincipalAmount = result.OSPrincipalAmount.Add(b.PrincipalAmount).Sub(b.PrincipalAmountPaid)
			result.OSInterestAmount = result.OSInterestAmount.Add(b.InterestAmount).Sub(b.InterestAmountPaid)
			result.TotalOSAmount = result.TotalOSAmount.Add(b.PrincipalAmount).Add(b.InterestAmount)
			// }
		}
	}

	return result, nil
}

// Delete implements Interface.
func (i *impl) Delete(ctx context.Context, params entity.DeleteLoanTransactionParams) (entity.LoanTransaction, error) {
	panic("unimplemented")
}

// Get implements Interface.
func (i *impl) Get(ctx context.Context, params entity.GetLoanTransactionParams) (entity.LoanTransaction, error) {
	panic("unimplemented")
}

// List implements Interface.
func (i *impl) List(ctx context.Context, params entity.ListLoanTransactionParams) ([]entity.LoanTransaction, entity.Pagination, error) {
	panic("unimplemented")
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateLoanTransactionParams) (entity.LoanTransaction, error) {
	panic("unimplemented")
}

// WithTx implements Interface.
func (i *impl) WithTx(ctx context.Context, tx *sql.Tx) Interface {
	panic("unimplemented")
}
