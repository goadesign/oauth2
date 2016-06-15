//************************************************************************//
// unnamed API: Application Media Types
//
// Generated with goagen v0.0.1, command line:
// $ goagen
// --design=github.com/goadesign/oauth2/design
// --out=$(GOPATH)/src/github.com/goadesign/oauth2
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package app

import "github.com/goadesign/goa"

// OAuth2ErrorMedia media type.
//
// Identifier: application/vnd.goa.example.oauth2.error+json
type OAuth2ErrorMedia struct {
	// Error returned by authorization server
	Error string `json:"error" xml:"error" form:"error"`
	// Human readable ASCII text providing additional information
	ErrorDescription *string `json:"error_description,omitempty" xml:"error_description,omitempty" form:"error_description,omitempty"`
	// A URI identifying a human-readable web page with information about the error
	ErrorURI *string `json:"error_uri,omitempty" xml:"error_uri,omitempty" form:"error_uri,omitempty"`
}

// Validate validates the OAuth2ErrorMedia media type instance.
func (mt *OAuth2ErrorMedia) Validate() (err error) {
	if mt.Error == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "error"))
	}

	if !(mt.Error == "invalid_request" || mt.Error == "invalid_client" || mt.Error == "invalid_grant" || mt.Error == "unauthorized_client" || mt.Error == "unsupported_grant_type") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.error`, mt.Error, []interface{}{"invalid_request", "invalid_client", "invalid_grant", "unauthorized_client", "unsupported_grant_type"}))
	}
	return
}

// TokenMedia media type.
//
// Identifier: application/vnd.goa.example.oauth2.token+json
type TokenMedia struct {
	// The access token issued by the authorization server
	AccessToken string `json:"access_token" xml:"access_token" form:"access_token"`
	// The lifetime in seconds of the access token
	ExpiresIn *int `json:"expires_in,omitempty" xml:"expires_in,omitempty" form:"expires_in,omitempty"`
	// The refresh token
	RefreshToken *string `json:"refresh_token,omitempty" xml:"refresh_token,omitempty" form:"refresh_token,omitempty"`
	// The scope of the access token
	Scope *string `json:"scope,omitempty" xml:"scope,omitempty" form:"scope,omitempty"`
	// The type of the token issued, e.g. "bearer" or "mac"
	TokenType string `json:"token_type" xml:"token_type" form:"token_type"`
}

// Validate validates the TokenMedia media type instance.
func (mt *TokenMedia) Validate() (err error) {
	if mt.AccessToken == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "access_token"))
	}
	if mt.TokenType == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "token_type"))
	}

	return
}
