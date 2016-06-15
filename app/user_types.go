//************************************************************************//
// unnamed API: Application User Types
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

// Payload sent by client to obtain refresh and access token or to refresh an access token.
// see https://tools.ietf.org/html/rfc6749#section-4.1.3 and https://tools.ietf.org/html/rfc6749#section-6
type tokenPayload struct {
	// The authorization code received from the authorization server, used for initial refresh and access token request
	Code *string `json:"code,omitempty" xml:"code,omitempty" form:"code,omitempty"`
	// Value MUST be set to "authorization_code" when obtaining initial refresh and access token.
	// Value MUST be set to "refresh_token" when refreshing an access token.
	GrantType *string `json:"grant_type,omitempty" xml:"grant_type,omitempty" form:"grant_type,omitempty"`
	// The redirect_uri parameter specified when making the authorize request to obtain the authorization code, used for initial refresh and access token request
	RedirectURI *string `json:"redirect_uri,omitempty" xml:"redirect_uri,omitempty" form:"redirect_uri,omitempty"`
	// The refresh token issued to the client, used for refreshing an access token
	RefreshToken *string `json:"refresh_token,omitempty" xml:"refresh_token,omitempty" form:"refresh_token,omitempty"`
	// The scope of the access request, used for refreshing an access token
	Scope *string `json:"scope,omitempty" xml:"scope,omitempty" form:"scope,omitempty"`
}

// Validate validates the tokenPayload type instance.
func (ut *tokenPayload) Validate() (err error) {
	if ut.GrantType == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "grant_type"))
	}

	if ut.GrantType != nil {
		if !(*ut.GrantType == "authorization_code" || *ut.GrantType == "refresh_token") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.grant_type`, *ut.GrantType, []interface{}{"authorization_code", "refresh_token"}))
		}
	}
	return
}

// Publicize creates TokenPayload from tokenPayload
func (ut *tokenPayload) Publicize() *TokenPayload {
	var pub TokenPayload
	if ut.Code != nil {
		pub.Code = ut.Code
	}
	if ut.GrantType != nil {
		pub.GrantType = *ut.GrantType
	}
	if ut.RedirectURI != nil {
		pub.RedirectURI = ut.RedirectURI
	}
	if ut.RefreshToken != nil {
		pub.RefreshToken = ut.RefreshToken
	}
	if ut.Scope != nil {
		pub.Scope = ut.Scope
	}
	return &pub
}

// Payload sent by client to obtain refresh and access token or to refresh an access token.
// see https://tools.ietf.org/html/rfc6749#section-4.1.3 and https://tools.ietf.org/html/rfc6749#section-6
type TokenPayload struct {
	// The authorization code received from the authorization server, used for initial refresh and access token request
	Code *string `json:"code,omitempty" xml:"code,omitempty" form:"code,omitempty"`
	// Value MUST be set to "authorization_code" when obtaining initial refresh and access token.
	// Value MUST be set to "refresh_token" when refreshing an access token.
	GrantType string `json:"grant_type" xml:"grant_type" form:"grant_type"`
	// The redirect_uri parameter specified when making the authorize request to obtain the authorization code, used for initial refresh and access token request
	RedirectURI *string `json:"redirect_uri,omitempty" xml:"redirect_uri,omitempty" form:"redirect_uri,omitempty"`
	// The refresh token issued to the client, used for refreshing an access token
	RefreshToken *string `json:"refresh_token,omitempty" xml:"refresh_token,omitempty" form:"refresh_token,omitempty"`
	// The scope of the access request, used for refreshing an access token
	Scope *string `json:"scope,omitempty" xml:"scope,omitempty" form:"scope,omitempty"`
}

// Validate validates the TokenPayload type instance.
func (ut *TokenPayload) Validate() (err error) {
	if ut.GrantType == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`response`, "grant_type"))
	}

	if !(ut.GrantType == "authorization_code" || ut.GrantType == "refresh_token") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError(`response.grant_type`, ut.GrantType, []interface{}{"authorization_code", "refresh_token"}))
	}
	return
}
