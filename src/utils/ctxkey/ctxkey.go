package ctxkey

type CtxKey string

const (
	USER_ID    = CtxKey("user_id")
	JWT_CLAIMS = CtxKey("jwt_claims")
)
