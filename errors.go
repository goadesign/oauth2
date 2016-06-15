package oauth2

import (
	"fmt"

	"github.com/goadesign/goa"
	"github.com/goadesign/oauth2/app"
)

type (
	// Error is the interface providing the implementation needed to return error
	// responses to OAuth2 requests as described in
	// https://tools.ietf.org/html/rfc6749#section-5.2
	Error interface {
		// Code is the error code and must have one of the values defines by the package
		// constants.
		Code() ErrorCode

		// Description is a human-readable ASCII text providing additional information, used
		// to assist the client developer in understanding the error that occurred.
		Description() string

		// URI identifying a human-readable web page with information about the error, used
		// to provide the client developer with additional information about the error.
		URI() string

		// Error implements the error interface.
		Error() string
	}

	// ErrorCode is the OAuth2 error code enum.
	ErrorCode string

	// oauth2Error is a simple implementation of Error.
	oauth2Error struct {
		ErrorCode        ErrorCode `json:"error"`
		ErrorDescription string    `json:"error_description,omitempty"`
		ErrorURI         string    `json:"error_uri,omitempty"`
	}
)

const (
	// ErrInvalidRequest is the error returned when the request is missing a required parameter,
	// includes an unsupported parameter value (other than grant type), repeats a parameter,
	// includes multiple credentials,  utilizes more than one mechanism for authenticating the
	// client, or is otherwise malformed.
	ErrInvalidRequest ErrorCode = "invalid_request"

	// ErrInvalidClient is the error returned when the client authentication failed (e.g., unknown
	// client, no client authentication included, or unsupported authentication method).
	ErrInvalidClient = "invalid_client"

	// ErrInvalidGrant is the error returned when the provided authorization grant (e.g.,
	// authorization code, resource owner credentials) or refresh token is invalid, expired,
	// revoked, does not match the redirection URI used in the authorization request, or was
	// issued to another client.
	ErrInvalidGrant = "invalid_grant"

	// ErrUnauthorizedClient is the error returned when the authenticated client is not authorized
	// to use this authorization grant type.
	ErrUnauthorizedClient = "unauthorized_client"

	// ErrUnsupportedGrantType is the error returned when the authorization grant type is not
	// supported by the authorization server.
	ErrUnsupportedGrantType = "unsupported_grant_type"

	// ErrInvalidScope is the error returned when the requested scope is invalid, unknown,
	// malformed, or exceeds the scope granted by the resource owner.
	ErrInvalidScope = "invalid_scope"
)

var (
	// ErrUnauthorized is the error returned for unauthorized requests.
	ErrUnauthorized = goa.NewErrorClass("unauthorized", 401)

	// MissingClientID is the response returned upon receiving a Authorize request with no
	// "client_id" query string.
	MissingClientID = errorToMedia(NewError(ErrInvalidRequest, "missing client ID", ""))

	// MissingCode is the response returned upon receiving a GetToken request with no
	// "code" form value in the request body.
	MissingCode = errorToMedia(NewError(ErrInvalidRequest, "missing authorization code", ""))

	// BadResponseType is the response returned upon receiving a Authorize request with a
	// response type set to something else than "code".
	BadResponseType = errorToMedia(NewError(ErrInvalidGrant, `only "code" response type is supported`, ""))

	// MissingRedirect is the response returned upon receiving a Authorize request with no
	// "redirect_uri" query string.
	MissingRedirect = errorToMedia(NewError(ErrInvalidRequest, "missing redirect URI", ""))

	// InvalidRedirect is the response returned upon receiving a Authorize request with a
	// malformed redirect URI.
	InvalidRedirect = errorToMedia(NewError(ErrInvalidRequest, "redirect URI must be a valid absolute URL", ""))

	// MalformedBody is the response returned upon receiving a GetToken request with a malformed
	// (non x-www-form-urlencoded) body.
	MalformedBody = errorToMedia(NewError(ErrInvalidRequest, "malformed body", ""))

	// InvalidGrantType is the response returned upon receiving a GetToken request with an
	// invalid grant_type form value.
	InvalidGrantType = errorToMedia(NewError(ErrInvalidGrant, `invalid grant type, must be "authorization_code" or "refresh_token"`, ""))

	// MissingRefreshToken is the response returned upon receiving a GetToken request with
	// grant type "refresh_token" and no refresh token.
	MissingRefreshToken = errorToMedia(NewError(ErrInvalidGrant, `grant type "refresh_token" requires a "refresh_token" value`, ""))
)

// NewError creates an error suitable to be returned in the body of OAuth2 error responses.
func NewError(code ErrorCode, description, uri string) Error {
	return &oauth2Error{code, description, uri}
}

// errorToMedia converts an error into a *app.OAuth2ErrorMedia. If e implements Error then the
// corresponding methods are used to build the content of the media struct otherwise a generic
// "invalid_request" response is returned.
func errorToMedia(e error) *app.OAuth2ErrorMedia {
	err, ok := e.(Error)
	if !ok {
		return &app.OAuth2ErrorMedia{Error: "invalid_request"}
	}
	m := &app.OAuth2ErrorMedia{Error: string(err.Code())}
	if d := err.Description(); d != "" {
		m.ErrorDescription = &d
	}
	if u := err.URI(); u != "" {
		m.ErrorURI = &u
	}
	return m
}

// oauth2Error implements Error.
func (e *oauth2Error) Code() ErrorCode     { return e.ErrorCode }
func (e *oauth2Error) Description() string { return e.ErrorDescription }
func (e *oauth2Error) URI() string         { return e.ErrorURI }
func (e *oauth2Error) Error() string       { return fmt.Sprintf("%v %s", e.ErrorCode, e.ErrorDescription) }
