package endpoint

import "github.com/josephspurrier/octane/example/app"

// Healthcheck -
// swagger:route GET / healthcheck HealthcheckGET
//
// Returns an OK message to show the application is functioning.
//
// You can use this endpoint as a healthcheck for a load balancer.
//
// Responses:
//   200: OKResponse
func Healthcheck(c *app.Context) error {
	return c.OKResponse("OK")
}
