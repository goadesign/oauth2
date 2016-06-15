package oauth2

import "golang.org/x/net/context"

// Private type used to key context.
type key int

// Context key value
const clientIDKey key = 1

// WithClientID creates a new context containing the given client ID that can be retrieved with
// ContextClientID.
func WithClientID(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, clientIDKey, clientID)
}

// ContextClientID extracts the client ID from the given context.
func ContextClientID(ctx context.Context) string {
	if cid := ctx.Value(clientIDKey); cid != nil {
		return cid.(string)
	}
	return ""
}
