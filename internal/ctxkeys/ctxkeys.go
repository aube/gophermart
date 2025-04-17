package ctxkeys

type contextKey string

const (
	UserID       contextKey = "userID"
	AuthTokenKey contextKey = "authToken"
	RequestIDKey contextKey = "requestID"
)
