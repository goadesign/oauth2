package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
	. "github.com/goadesign/oauth2/design/public"
)

// OAuth2 initializes the design definitions needed to implement a OAuth2 provider.
// This function defines the OAuth2Provider resource which is implemented by the
// OAuth2ProviderController defined in the parent package. This controller implements the
// authorization and token endpoints.
//
// The value returned can be used in the design of services to secure access to resources using
// OAuth2.  The scheme uses the standard "Authorization Code Grant" flow described in RFC 6749 where
// the client application must first retrieve an authorization code by requesting access to the
// resource owner then use the code to request a refresh and access token.
//
// authorizationEndpoint is the request path to the authorization endpoint as described by
// https://tools.ietf.org/html/rfc6749#section-3.1. This endpoint receives requests from the
// resource owner to grant access to the client.  It responds with a redirect to a preconfigured URI
// and provides the authorization code as a query string.
//
// tokenEndpoint is the request path to the token endpoint as described by
// https://tools.ietf.org/html/rfc6749#section-3.2. This endpoint exchanges authorization codes for
// refresh and access tokens and also refreshes access tokens given a refresh token.
//
// dsl is an optional anonymous function that may define scopes for the grants associated with the
// access tokens.
//
// Example:
//
//    var OAuth2Sec = OAuth2("/oauth2/auth", "/oauth2/token", func() {
//        Scope("api:read", "Scope granting read access")
//        Scope("api:write", "Scope granting write access")
//    })
//
func OAuth2(authorizationEndpoint, tokenEndpoint string, dsl ...func()) *SecuritySchemeDefinition {
	// The resource that implements the OAuth2 standard defined by RFC 6749.
	// See https://tools.ietf.org/html/rfc6749
	var _ = Resource("oauth2_provider", func() {
		Description("This resource implements the OAuth2 authorization code flow")

		Action("authorize", func() {
			Description("Authorize OAuth2 client")
			Routing(GET(authorizationEndpoint))
			Params(func() {
				Param("response_type", String, `Value MUST be set to "code"`, func() {
					Enum("code")
				})
				Param("client_id", String, "The client identifier")
				Param("redirect_uri", String, "Redirection endpoint")
				Param("scope", String, "The scope of the access request")
				Param("state", String, "An opaque value used by the client to maintain state between the request and callback")
				Required("response_type", "client_id")
			})
			Response(Found, func() {
				Headers(func() {
					Header("Location", String, "Redirect URL containing the authorization code and state param if any")
				})
			})
			Response(BadRequest, OAuth2ErrorMedia)
		})

		Action("get_token", func() {
			Description("Get access token from authorization code or refresh token")
			Routing(POST(tokenEndpoint))
			Security(OAuth2ClientBasicAuth)
			Payload(OAuth2TokenPayload)
			Response(OK, OAuth2TokenMedia)
			Response(BadRequest, OAuth2ErrorMedia)
		})
	})

	// Define security scheme
	return OAuth2Security("OAuth2", func() {
		// AccessCodeFlow defines a "Authorization Code" OAuth2 flow
		// see https://tools.ietf.org/html/rfc6749#section-1.3.1
		AccessCodeFlow(authorizationEndpoint, tokenEndpoint)

		// Run the DSL which sets up optional scopes.
		if len(dsl) > 0 {
			dsl[0]()
		}
	})

}

// OAuth2ClientBasicAuth defines the basic auth used to make requests to the token endpoint.  The
// username and password must correspond to the client id and secret and be encoded using form
// encoding as described in https://tools.ietf.org/html/rfc6749#section-2.3.1
var OAuth2ClientBasicAuth = BasicAuthSecurity("oauth2_client_basic_auth", func() {
	Description("Basic auth used by client to make the requests needed to retrieve and refresh access tokens")
})
