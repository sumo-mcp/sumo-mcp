package api

import (
	"net/http"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/sumo-mcp/sumoapi-go"
)

// API defines the operations available in the Sumo API.
type API = sumoapi.Client

// wrapper is a wrapper around API to allow for custom
// behavior, e.g. there are bugs that can be fixed
// on the wrapper side.
type wrapper struct{ API }

func New(httpClient *http.Client) API {
	return wrapper{sumoapi.New(sumoapi.WithHTTPClient(httpClient))}
}

// TypeSchemas returns the registered type schemas.
func TypeSchemas() map[reflect.Type]*jsonschema.Schema {
	return sumoapi.TypeSchemas()
}
