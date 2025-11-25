module github.com/sumo-mcp/sumo-mcp

go 1.25.0

require (
	github.com/google/jsonschema-go v0.3.0
	github.com/modelcontextprotocol/go-sdk v1.1.0
	github.com/sumo-mcp/sumoapi-go v0.5.0
)

replace github.com/modelcontextprotocol/go-sdk => github.com/matheuscscp/mcp-go-sdk v0.0.0-20251125234243-c1ad726cdaba

require (
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
)
