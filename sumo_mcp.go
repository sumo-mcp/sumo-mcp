package sumomcp

import (
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/sumo-mcp/sumo-mcp/internal/api"
	"github.com/sumo-mcp/sumo-mcp/internal/server"
)

// webOrigin is the origin of the static website (https://sumo-mcp.com) that is
// allowed to call the MCP endpoint directly from the browser. It must be a
// distinct origin from the API itself (https://api.sumo-mcp.com), so it is both
// registered as a trusted origin for the SDK's cross-origin protection and
// reflected back in CORS headers.
const webOrigin = "https://sumo-mcp.com"

func init() {
	a := api.New(&http.Client{Timeout: 5 * time.Second})
	s := server.New("http", a)

	cop := http.NewCrossOriginProtection()
	if err := cop.AddTrustedOrigin(webOrigin); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/mcp", withCORS(mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return s },
		&mcp.StreamableHTTPOptions{CrossOriginProtection: cop},
	)))

	functions.HTTP("sumo-mcp", mux.ServeHTTP)
}

// withCORS allows the static website to call the MCP endpoint from the browser.
// It answers CORS preflight (OPTIONS) requests and adds the response headers
// browsers require to read the streamable-HTTP session id.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == webOrigin {
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", webOrigin)
			h.Add("Vary", "Origin")
			h.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Content-Type, Mcp-Session-Id, Mcp-Protocol-Version")
			h.Set("Access-Control-Expose-Headers", "Mcp-Session-Id")
			h.Set("Access-Control-Max-Age", "86400")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
