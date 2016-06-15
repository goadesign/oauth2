//go:generate goagen app -d github.com/goadesign/oauth2/design

/*
Package oauth2 contains the implementation of the OAuth2 provider controller.
*/
package oauth2

import (
	"net/http"
	"net/url"

	"github.com/goadesign/goa"
	"github.com/goadesign/oauth2/app"

	"golang.org/x/net/context"
)

type (
	// ProviderController implements the OAuth2Provider resource.
	ProviderController struct {
		*goa.Controller
		provider Provider // User provided implementation
	}

	// Provider is the interface that provides the actual implementation for the authorize
	// and token endpoints.
	Provider interface {
		// Authorize implements https://tools.ietf.org/html/rfc6749#section-4.1.1
		// Given a client identifier the implementation must verify that the redirect URI
		// matches the pre-registered URI. The implementation should also validate the
		// scope. Upon success Authorize should return the authorization code and a nil
		// error. Upon failure the error should implement Error otherwise a generic error
		// HTTP response is sent back to the client.
		Authorize(clientID, scope, redirectURI string) (code string, err error)

		// Exchange implements https://tools.ietf.org/html/rfc6749#section-4.1.3 It must
		// check that the given authorization code was generated for the client with the
		// given identifier and that the redirect URI matches the pre-registered redirect
		// URI. Upon success it should return a refresh and access token pair as well as an
		// optional expiration deadline in seconds. Upon failure the error should implement
		// Error otherwise a generic error HTTP response is sent back to the client.
		Exchange(clientID, code, redirectURI string) (refreshToken, accessToken string, expiresIn int, err error)

		// Refresh implements https://tools.ietf.org/html/rfc6749#section-6
		// It must check that the given refresh token and scope are valid.  Upon success it
		// should return a valid access token and optionally a new refresh token and
		// expiration deadline in seconds. Upon failure the error should implement Error
		// otherwise a generic error HTTP response is sent back to the client.
		Refresh(refreshToken, scope string) (newRefreshToken, accessToken string, expiresIn int, err error)

		// Authenticate performs client authentication as described in
		// https://tools.ietf.org/html/rfc6749#section-2.3
		// It should return nil if the client is authorized, a non-nil error otherwise.
		// The error message is returned in the Unauthorized response body.
		Authenticate(clientID, clientSecret string) error
	}
)

// NewProviderController creates a OAuth2Provider controller.
func NewProviderController(service *goa.Service, provider Provider) *ProviderController {
	return &ProviderController{
		Controller: service.NewController("OAuth2ProviderController"),
		provider:   provider,
	}
}

// NewOAuth2ClientBasicAuthMiddleware creates the security middleware to be used for authenticating
// the client GetToken requests. The given callback must validate the client credentials.
func NewOAuth2ClientBasicAuthMiddleware(provider Provider) goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			// Retrieve basic auth info, TBD these are urlencoded as per the spec...
			clientID, clientSecret, ok := req.BasicAuth()
			if !ok {
				return ErrUnauthorized("missing auth")
			}

			// Validate creds
			if err := provider.Authenticate(clientID, clientSecret); err != nil {
				return ErrUnauthorized(err)
			}

			// Store client ID in context and proceed
			ctx = WithClientID(ctx, clientID)
			return h(ctx, rw, req)
		}
	}
}

// Authorize is a request made by the resource owner to grant access to the client.  It redirects
// to the client using a pre-registered redirect URI. The redirect URL contains the authorization
// code as a query string value.
func (c *ProviderController) Authorize(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	var (
		query        = req.URL.Query()
		clientID     = query.Get("client_id")
		responseType = query.Get("response_type")
		redirectURI  = query.Get("redirect_uri")
		scope        = query.Get("scope")
		state        = query.Get("state")
	)
	// Ensure there is a client identifier
	if clientID == "" {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingClientID)
	}

	// Validate grant type
	if responseType != "code" {
		return c.Service.Send(ctx, http.StatusBadRequest, BadResponseType)
	}

	// Ensure there is a redirect URI
	if redirectURI == "" {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingRedirect)
	}

	// Validate redirect URI
	u, err := url.Parse(redirectURI)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return c.Service.Send(ctx, http.StatusBadRequest, InvalidRedirect)
	}

	// Retrieve auth code
	code, err := c.provider.Authorize(clientID, scope, redirectURI)
	if err != nil {
		return c.Service.Send(ctx, http.StatusBadRequest, errorToMedia(err))
	}

	// Write code to URL query
	q := u.Query()
	q.Set("code", code)

	// Keep state originally provided by client if any
	if state != "" {
		q.Set("state", state)
	}

	u.RawQuery = q.Encode()
	rw.Header().Set("Location", u.String())

	return c.Service.Send(ctx, http.StatusFound, nil)
}

// GetToken runs the get_token action.
func (c *ProviderController) GetToken(ctx context.Context, rw http.ResponseWriter, grantType string,
	code, redirectURI, refreshToken, scope *string) error {

	if grantType == "authorization_code" {
		return c.exchange(ctx, rw, code, redirectURI)
	}
	if grantType == "refresh_token" {
		return c.refresh(ctx, rw, refreshToken, scope)
	}
	return c.Service.Send(ctx, http.StatusBadRequest, InvalidGrantType)
}

// exchange returns a pair of refresh and access tokens from an authorization code.
func (c *ProviderController) exchange(ctx context.Context, rw http.ResponseWriter, code, redirectURI *string) error {
	// Ensure there is a client identifier
	clientID := ContextClientID(ctx)
	if clientID == "" {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingClientID)
	}

	// Ensure there is a code
	if code == nil {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingCode)
	}

	// Ensure there is a redirect URI
	if redirectURI == nil {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingRedirect)
	}

	// Validate redirect URI
	u, err := url.Parse(*redirectURI)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return c.Service.Send(ctx, http.StatusBadRequest, InvalidRedirect)
	}

	// Retrieve tokens code
	refreshToken, accessToken, expiresIn, err := c.provider.Exchange(clientID, *code, *redirectURI)
	if err != nil {
		return c.Service.Send(ctx, http.StatusBadRequest, errorToMedia(err))
	}

	m := app.TokenMedia{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}
	if refreshToken != "" {
		m.RefreshToken = &refreshToken
	}
	if expiresIn != 0 {
		m.ExpiresIn = &expiresIn
	}

	rw.Header().Set("Content-Type", "application/json")

	return c.Service.Send(ctx, http.StatusOK, &m)
}

// refresh refreshes an access token given a refresh token.
func (c *ProviderController) refresh(ctx context.Context, rw http.ResponseWriter, refreshToken, scope *string) error {
	// Ensure there is a refresh token
	if refreshToken == nil {
		return c.Service.Send(ctx, http.StatusBadRequest, MissingRefreshToken)
	}

	// Retrieve tokens
	var s string
	if scope != nil {
		s = *scope
	}
	rToken, aToken, expiresIn, err := c.provider.Refresh(*refreshToken, s)
	if err != nil {
		return c.Service.Send(ctx, http.StatusBadRequest, errorToMedia(err))
	}

	m := app.TokenMedia{
		AccessToken: aToken,
		TokenType:   "Bearer",
	}
	if rToken != "" {
		m.RefreshToken = &rToken
	}
	if expiresIn != 0 {
		m.ExpiresIn = &expiresIn
	}
	if scope != nil {
		m.Scope = scope
	}

	rw.Header().Set("Content-Type", "application/json")

	return c.Service.Send(ctx, http.StatusOK, &m)
}
