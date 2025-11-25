package api

import (
	"net/http"

	"github.com/sumo-mcp/sumoapi-go"
)

// client is a wrapper around sumoapi.Client to allow
// for custom behavior, e.g. there are bugs that can
// be fixed on the client side.
type client struct{ sumoapi.Client }

func New(httpClient *http.Client) sumoapi.Client {
	return client{sumoapi.New(sumoapi.WithHTTPClient(httpClient))}
}
