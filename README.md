# OAuth2 Provider

This repo provides both the design and implementation for a complete OAuth2 provider. This provider
implements the "Authorization Code" flow as described by
[RFC 6749 section 1.3.1](https://tools.ietf.org/html/rfc6749#section-1.3.1).

Looking at the [Protocol Flow](https://tools.ietf.org/html/rfc6749#section-1.2) this repo implements
the `Authorization Server` (and your service the `Resource Server`). Included is the design and
implementation for three different types of requests:

* The request coming from the resource owner that authorizes the client. As a concrete example if
  your service implemented Facebook then the request would be the post sent by the end user upon
  accepting the "Request for Permission" displayed when adding an app such as Spotify. The result of
  this request is a redirection to a pre-configured Spotify URI with an authorization code in the
  redirect URL query string. The Spotify app can now use the authorization code to retrieve a
  refresh and access token from Facebook.

* The request coming from the client to exchange the authorization code obtained in the step above
  for a refresh and access tokens. Keeping with the Facebook example this request comes from Spotify
  after it receives the authorization code from the redirect to retrieve a refresh and access token.
  At this point Spotify can make requests on behalf of the end user using the access token to
  authenticate against our service (Facebook).

* The request coming from the client to renew an access token using a refresh token obtained in the
  step above. In the Facebook example this request would be again from Spotify when or before an
  access token expires to retrieve a fresh new access token.

The last two requests are depicted in "Figure 2: Refreshing an Expired Access Token" in
[section 1.5 of the RFC](https://tools.ietf.org/html/rfc6749#section-1.5) as steps `(A)`, `(B)` and
`(G)`, `(H)`.

The complete flow is depicted in "Figure 3: Authorization Code Flow" in
[section 4.1 of the RFC](https://tools.ietf.org/html/rfc6749#section-4.1). Remember that this repo
implements the "Authorization Server" and that the "Client" is the 3rd party service wanting to use
OAuth2 to send requests to your service on behalf of the "Resource Owner" (i.e. end user of your
service).

## Usage

Using the provider consists of two parts:

1. Import the design package in the design of the API that implements the OAuth2 provider and call
   the `OAuth2` function which creates the security scheme. Use the scheme to secure the API
   endpoints as needed. See [Design](#design) below.
2. Implement and mount the oauth2 provider controller onto the service.
   See [Implement](#implement) below.

### Design

First import the OAuth2 design package in your design:

```go
package design

import (
    . "github.com/goadesign/design"
    . "github.com/goadesign/design/apidsl"
    . "github.com/goadesign/oauth2"
)
```

Then call the `OAuth2` function to create the OAuth2 security scheme. This function accepts three
arguments:

* `authorizationEndpoint` is the request path to the authorization endpoint as described by
https://tools.ietf.org/html/rfc6749#section-3.1. This endpoint receives requests from the
resource owner to grant access to the client.  It responds with a redirect to a preconfigured URI
and provides the authorization code as a query string.

* `tokenEndpoint` is the request path to the token endpoint as described by
https://tools.ietf.org/html/rfc6749#section-3.2. This endpoint exchanges authorization codes and
refresh tokens for access tokens.

* `dsl` is an optional anonymous function that may define scopes for the grants associated with the
access tokens.

Here is an example creating a OAuth2 security scheme using `/oauth2/authorize` and `/oauth2/token`
as endpoints and definiting two scopes `api:read` and `api:write`:

```go
var OAuth2Sec = OAuth2("/oauth2/authorize", "/oauth2/token", func() {
    Scope("api:read")
    Scope("api:write")
})
```

Not too bad uh? this scheme can now be used to secure endpoints in the service, for example:

```go
var _ = Resource("secure_resource", func() {
    Description("The actions of this resource require a valid OAuth2 access token with a 'api:read' scope")

    Security(OAuth2Sec, func() {
        Scope("api:read")
    })

    Action("super_secure", func() {
        Description("This action also requires the 'api:write' scope")

        Security(OAuth2Sec, func() {
            Scope("api:write")
        })

        // ...
    })

    // ...
})
```

That's it! You now have designed an OAuth2 enabled service. Onto implementation.

## Implement

Implementing the OAuth2 provider controller defined in the generated `oauth2_provider.go` is done by
instantiating the `ProviderController` struct provided in this package and using its methods to
implement the generated controller actions.

There are two actions genererated: `Authorize` and `GetToken`. The `ProviderController` struct
exposes methods with the same names that provide the implementation for the respective actions.

Concretly here is how the generated controller could be implemented. The following should be written
in the generated `oauth2_provider.go` file and replace the correspodning placeholder functions:

```go
// NewOAuth2ProviderController creates a OAuth2Provider controller.
func NewOAuth2ProviderController(service *goa.Service, provider oauth2.Provider) *OAuth2ProviderController {
	return &OAuth2ProviderController{
		ProviderController: oauth2.NewProviderController(service, provider),
	}
}

// Authorize runs the authorize action.
func (c *OAuth2ProviderController) Authorize(ctx *app.AuthorizeOauth2ProviderContext) error {
	return c.ProviderController.Authorize(ctx.Context, ctx.ResponseWriter, ctx.Request)
}

// GetToken runs the get_token action.
func (c *OAuth2ProviderController) GetToken(ctx *app.GetTokenOauth2ProviderContext) error {
	p := ctx.Payload
	return c.ProviderController.GetToken(ctx.Context, ctx.ResponseWriter, p.GrantType,
		p.Code, p.RedirectURI, p.RefreshToken, p.Scope)
}
```
 
As you can see instantiating the provider controller requires passing in a `provider` parameter of
type `oauth2.Provider`. This interface exposes four high level methods that makes it possible to
inject the actual OAuth2 authorization logic. The definition of the interface is:

```go
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
```

The bulk of the work thus consists of implementing this interface. The package then takes care of
validating the incoming requests, invoking the methods above in the right places and returning
properly formatted success or error responses.

The implementation of these methods can take advantage of the errors defined in this package. In
particular the `NewError` function should be used to create the instances of errors returned by the
methods.

The [security example](https://github.com/goadesign/examples/blob/master/security) contains a complete
implementation of a OAuth2 provider as well as instructions for how to use the generated client to
make requests to go through the authorization flow.
