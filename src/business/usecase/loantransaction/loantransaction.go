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
	"github.com/irdaislakhuafa/go-sdk/convert"
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
		Pay(ctx context.Context, params entity.PayLoanTransactionParams) (entity.LoanTransaction, error)
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

	// get setting
	var err error
	var EOD_DATE time.Time
	{
		setting, err := i.dom.Setting.Get(ctx, entity.GetSettingParams{
			Name:      entity.SETTING_NAME_EOD_DATE,
			IsDeleted: 0,
		})
		if err != nil {
			code := errors.GetCode(err)
			if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
				return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_NOT_FOUND)
			}
			return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}

		EOD_DATE, err = time.Parse(time.DateOnly, setting.Value)
		if err != nil {
			return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_EOD_DATE_INVALID)
		}
	}

	// get user
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
	result.InvoiceNumber = result.GenInvoiceNumber(result.ID, result.UserID)

	_, err = dLoanTransaction.Update(ctx, entity.UpdateLoanTransactionParams{
		InvoiceNumber: result.InvoiceNumber,
		Notes:         result.Notes,
		UserID:        result.UserID,
		User:          result.User,
		LoanID:        result.LoanID,
		Loan:          result.Loan,
		Amount:        result.Amount.InexactFloat64(),
		ID:            result.ID,
	})
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	// generate loan billing
	listBilling := []entity.LoanBilling{}
	switch loan.RepaymentType {
	case entity.REPAYMENT_TYPE_WEEKS:
		billDate := EOD_DATE
		for i := 0; i < int(loan.RepaymentDuration); i++ {
			billDate = billDate.AddDate(0, 0, 7)
			listBilling = append(listBilling, entity.LoanBilling{
				LoanTransactionID: result.ID,
				BillDate:          billDate,
			})
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
	// if err := tx.Commit(); err != nil {
	// 	return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	// }

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

	// get setting
	setting, err := i.dom.Setting.Get(ctx, entity.GetSettingParams{
		Name:      entity.SETTING_NAME_EOD_DATE,
		IsDeleted: 0,
	})
	if err != nil {
		return result, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	EOD_DATE, err := time.Parse(time.DateOnly, setting.Value)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_EOD_DATE_INVALID)
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
		now := EOD_DATE
		for _, b := range billings {
			// get current billing
			isCurrent := b.BillDate.Before(now) && (b.PrincipalAmountPaid.LessThan(b.PrincipalAmount) || b.InterestAmountPaid.LessThan(b.InterestAmount))
			if isCurrent {
				result.CurrentBillDate = &b.BillDate
				result.ListBilledBilling = append(result.ListBilledBilling, b)
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

func (i *impl) Pay(ctx context.Context, params entity.PayLoanTransactionParams) (entity.LoanTransaction, error) {
	// validate params
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	// get setting
	setting, err := i.dom.Setting.Get(ctx, entity.GetSettingParams{
		Name:      entity.SETTING_NAME_EOD_DATE,
		IsDeleted: 0,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_NOT_FOUND)
		}
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	EOD_DATE, err := time.Parse(time.DateOnly, setting.Value)
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.SETTING_EOD_DATE_INVALID)
	}

	// get loan transaction
	lt, err := i.dom.LoanTransaction.Get(ctx, entity.GetLoanTransactionParams{
		ID: params.LoanTransactionID,
	})
	if err != nil {
		code := errors.GetCode(err)
		if code.IsOneOf(codes.CodeSQLRecordDoesNotExist) {
			return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_TRANSACTION_NOT_FOUND)
		}
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	// get billing
	billings, _, err := i.dom.LoanBilling.List(ctx, entity.ListLoanBillingParams{
		LoanTransactionID: params.LoanTransactionID,
		PaginationParams: entity.PaginationParams{
			Limit:     9_999_999,
			Page:      0,
			OrderBy:   "bill_date",
			OrderType: "desc",
		},
		IsDeleted:              0,
		UserID:                 params.UserID,
		BillDateLTE:            EOD_DATE,
		PrincipalAmountPaidLTE: convert.ToPointer(decimal.NewFromInt(0)),
		InterestAmountPaidLTE:  convert.ToPointer(decimal.NewFromInt(0)),
	})
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	if len(billings) == 0 {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.LOAN_BILLING_NOT_FOUND)
	}

	tx, err := i.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	dLoanBilling := i.dom.LoanBilling.WithTx(ctx, tx)

	// calculate amount
	amount := decimal.NewFromFloat(params.Amount)
	for i := range billings {
		b := billings[i]
		if amount.LessThanOrEqual(decimal.NewFromInt(0)) {
			break
		}
		amount = amount.Sub(b.PrincipalAmount).Sub(b.InterestAmount)
		b.PrincipalAmountPaid = b.PrincipalAmount
		b.InterestAmountPaid = b.InterestAmount
		_, err := dLoanBilling.Update(ctx, entity.UpdateLoanBillingParams{
			LoanTransactionID:   b.LoanTransactionID,
			BillDate:            b.BillDate,
			PrincipalAmount:     b.PrincipalAmount,
			PrincipalAmountPaid: b.PrincipalAmountPaid,
			InterestAmount:      b.InterestAmount,
			InterestAmountPaid:  b.InterestAmountPaid,
			ID:                  b.ID,
		})
		if err != nil {
			return entity.LoanTransaction{}, errors.NewWithCode(errors.GetCode(err), err.Error())
		}
		billings[i] = b
	}

	// if amount not match then return error
	if !amount.Equal(decimal.NewFromInt(0)) {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeBadRequest, errmessages.TRANSACTION_AMOUNT_NOT_MATCH)
	}

	// commit tx
	if err := tx.Commit(); err != nil {
		return entity.LoanTransaction{}, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	lt.LoanBilling = billings
	return lt, nil

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
	if err := i.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return []entity.LoanTransaction{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	results, pagination, err := i.dom.LoanTransaction.List(ctx, params)
	if err != nil {
		return []entity.LoanTransaction{}, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), err.Error())
	}

	return results, pagination, nil
}

// Update implements Interface.
func (i *impl) Update(ctx context.Context, params entity.UpdateLoanTransactionParams) (entity.LoanTransaction, error) {
	panic("unimplemented")
}
