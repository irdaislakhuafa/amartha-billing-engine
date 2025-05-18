package errmessages

const (
	LOAN_NOT_FOUND = "loan not found"

	USER_NOT_REGISTERED     = "user not registered"
	USER_ALREADY_REGISTERED = "user already registered"
	USER_SUCCESS_REGISTERED = "user register successfully"
	USER_PASSWORD_NOT_MATCH = "password not match"
	USER_IS_DELINQUENT      = "this user is delinquent and cannot process loan transaction!"

	AUTH_HEADER_REQUIRED = "auth header is required to access this resource!"

	LOAN_REPAYMENT_TYPE_INVALID       = "invalid loan repayment type!"
	LOAN_REPAYMENT_TYPE_NOT_AVAILABLE = "loan repayment type not available!"
	LOAN_TRANSACTION_NOT_FOUND        = "loan transaction not found"

	SETTING_NOT_FOUND                              = "setting not found"
	SETTING_EOD_DATE_NOT_FOUND                     = "setting eod date not found"
	SETTING_EOD_DATE_INVALID                       = "setting eod date invalid"
	SETTING_LIMIT_BILLING_FOR_DELINQUENT_INVALID   = "setting limit billing for delinquent invalid"
	SETTING_LIMIT_BILLING_FOR_DELINQUENT_NOT_FOUND = "setting limit billing for delinquent not found"
	SETTING_ALREADY_EXISTS                         = "setting already exists"

	LOAN_BILLING_NOT_FOUND = "loan billing not found"

	TRANSACTION_AMOUNT_NOT_MATCH = "the amount paid does not match, amount can't be less than or greater than the outstanding amount"
	TRANSACTION_LIMIT            = "please complete your previous loan transaction billing before creating a new one"
)
