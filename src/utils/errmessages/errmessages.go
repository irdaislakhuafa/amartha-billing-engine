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
)
