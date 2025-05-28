package contextkey

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	UUIDKey   contextKey = "uuid"
	RefKey    contextKey = "ref"
	RefIDKey  contextKey = "ref_id"
)
