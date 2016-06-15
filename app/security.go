//************************************************************************//
// unnamed API: Application Security
//
// Generated with goagen v0.0.1, command line:
// $ goagen
// --design=github.com/goadesign/oauth2/design
// --out=$(GOPATH)/src/github.com/goadesign/oauth2
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import (
	"github.com/goadesign/goa"
	"golang.org/x/net/context"
	"net/http"
)

type (
	// Private type used to store auth handler info in request context
	authMiddlewareKey string
)

// UseOauth2ClientBasicAuthMiddleware mounts the oauth2_client_basic_auth auth middleware onto the service.
func UseOauth2ClientBasicAuthMiddleware(service *goa.Service, middleware goa.Middleware) {
	service.Context = context.WithValue(service.Context, authMiddlewareKey("oauth2_client_basic_auth"), middleware)
}

// NewOauth2ClientBasicAuthSecurity creates a oauth2_client_basic_auth security definition.
func NewOauth2ClientBasicAuthSecurity() *goa.BasicAuthSecurity {
	def := goa.BasicAuthSecurity{}
	def.Description = "Basic auth used by client to make the requests needed to retrieve and refresh access tokens"
	return &def
}

// handleSecurity creates a handler that runs the auth middleware for the security scheme.
func handleSecurity(schemeName string, h goa.Handler, scopes ...string) goa.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		scheme := ctx.Value(authMiddlewareKey(schemeName))
		am, ok := scheme.(goa.Middleware)
		if !ok {
			return goa.NoAuthMiddleware(schemeName)
		}
		ctx = goa.WithRequiredScopes(ctx, scopes)
		return am(h)(ctx, rw, req)
	}
}
