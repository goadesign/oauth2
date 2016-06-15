package public

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// OAuth2TokenMedia describes the response sent in case of successful access token request.
// See https://tools.ietf.org/html/rfc6749#section-5.1
var OAuth2TokenMedia = MediaType("application/vnd.goa.example.oauth2.token+json", func() {
	Description("OAuth2 access token request successful response, see https://tools.ietf.org/html/rfc6749#section-5.1")
	TypeName("TokenMedia")
	Attributes(func() {
		Attribute("access_token", String, "The access token issued by the authorization server")
		Attribute("token_type", String, `The type of the token issued, e.g. "bearer" or "mac"`)
		Attribute("expires_in", Integer, "The lifetime in seconds of the access token")
		Attribute("refresh_token", String, "The refresh token")
		Attribute("scope", String, "The scope of the access token")
		Required("access_token", "token_type")
	})
	View("default", func() {
		Attribute("access_token")
		Attribute("token_type")
		Attribute("expires_in")
		Attribute("refresh_token")
		Attribute("scope")
	})
})

// OAuth2ErrorMedia describes responses sent in case of invalid request to the provider endpoints.
// See https://tools.ietf.org/html/rfc6749#section-4.1.2.1
var OAuth2ErrorMedia = MediaType("application/vnd.goa.example.oauth2.error+json", func() {
	Description("OAuth2 error response, see https://tools.ietf.org/html/rfc6749#section-5.2")
	TypeName("OAuth2ErrorMedia")
	Attributes(func() {
		Attribute("error", String, "Error returned by authorization server", func() {
			Enum("invalid_request", "invalid_client", "invalid_grant", "unauthorized_client", "unsupported_grant_type")
		})
		Attribute("error_description", String, "Human readable ASCII text providing additional information")
		Attribute("error_uri", String, "A URI identifying a human-readable web page with information about the error")
		Required("error")
	})
	View("default", func() {
		Attribute("error")
		Attribute("error_description")
		Attribute("error_uri")
	})
})
