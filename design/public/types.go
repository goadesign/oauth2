package public

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// OAuth2TokenPayload describes the body sent by the client to obtain refresh and access tokens
// using an authorization code.
// See https://tools.ietf.org/html/rfc6749#section-4.1.3
//
// It also describes the body sent by the client to refresh a token.
// See https://tools.ietf.org/html/rfc6749#section-6
//
var OAuth2TokenPayload = Type("TokenPayload", func() {
	Description(`Payload sent by client to obtain refresh and access token or to refresh an access token.
see https://tools.ietf.org/html/rfc6749#section-4.1.3 and https://tools.ietf.org/html/rfc6749#section-6`)
	Attribute("grant_type", String, `Value MUST be set to "authorization_code" when obtaining initial refresh and access token.
Value MUST be set to "refresh_token" when refreshing an access token.`, func() {
		Enum("authorization_code", "refresh_token")
	})

	// Initial refresh and access token request payload
	Attribute("code", String, "The authorization code received from the authorization server, used for initial refresh and access token request")
	Attribute("redirect_uri", String, "The redirect_uri parameter specified when making the authorize request to obtain the authorization code, used for initial refresh and access token request")

	// Refresh token payload
	Attribute("refresh_token", String, "The refresh token issued to the client, used for refreshing an access token")
	Attribute("scope", String, "The scope of the access request, used for refreshing an access token")

	Required("grant_type")
})
