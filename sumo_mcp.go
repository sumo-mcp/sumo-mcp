package sumomcp

import (
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/sumo-mcp/sumo-mcp/internal/api"
	"github.com/sumo-mcp/sumo-mcp/internal/server"
)

func init() {
	a := api.New(&http.Client{Timeout: 5 * time.Second})
	s := server.New("http", a)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return s },
		&mcp.StreamableHTTPOptions{},
	))

	functions.HTTP("sumo-mcp", mux.ServeHTTP)
}
