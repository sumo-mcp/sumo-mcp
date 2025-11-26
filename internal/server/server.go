package server

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/sumo-mcp/sumo-mcp/internal/api"
)

//go:embed instructions.md
var instructions string

func New(version string, a api.API) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "sumo-mcp",
		Title:   "Sumo MCP Server",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: instructions,
		HasTools:     true,
	})

	addObjectTool(s, "search_rikishi",
		"Search for rikishi (sumo wrestlers).",
		a.SearchRikishi)

	addObjectTool(s, "get_rikishi",
		"Get detailed profile information about a specific rikishi (sumo wrestler).",
		a.GetRikishi)

	addObjectTool(s, "get_rikishi_stats",
		"Get overall performance stats about a specific rikishi (sumo wrestler).",
		a.GetRikishiStats)

	addObjectTool(s, "list_rikishi_matches",
		"List matches for a specific rikishi (sumo wrestler).",
		a.ListRikishiMatches)

	addObjectTool(s, "list_rikishi_matches_against_opponent",
		"List matches for a specific rikishi (sumo wrestler) against a specific opponent.",
		a.ListRikishiMatchesAgainstOpponent)

	addObjectTool(s, "get_basho",
		"Get detailed information about a specific basho (sumo tournament).",
		a.GetBasho)

	addObjectTool(s, "get_banzuke",
		"Get the banzuke (ranking list) for a specific basho (sumo tournament) division.",
		a.GetBanzuke)

	addObjectTool(s, "get_basho_with_torikumi",
		"Get detailed information about a specific basho (sumo tournament) along with the torikumi (bout schedule) for a specific day and division.",
		a.GetBashoWithTorikumi)

	addObjectTool(s, "list_kimarite",
		"List kimarite (winning techniques) along with their statistics.",
		a.ListKimarite)

	addObjectTool(s, "list_kimarite_matches",
		"List matches won by a specific kimarite (winning technique).",
		a.ListKimariteMatches)

	addObjectListTool(s, "list_measurement_changes",
		"List measurement changes for rikishi (sumo wrestlers).",
		a.ListMeasurementChanges)

	addObjectListTool(s, "list_rank_changes",
		"List rank changes for rikishi (sumo wrestlers).",
		a.ListRankChanges)

	addObjectListTool(s, "list_shikona_changes",
		"List shikona (ring name) changes for rikishi (sumo wrestlers).",
		a.ListShikonaChanges)

	return s
}

func addObjectTool[In, Out any](s *mcp.Server, name, desc string, fn func(context.Context, In) (*Out, error)) {
	schemaOpts := &jsonschema.ForOptions{
		TypeSchemas: api.TypeSchemas(),
	}

	inputSchema, err := jsonschema.For[In](schemaOpts)
	if err != nil {
		panic(fmt.Sprintf("error inferring input schema for %s: %v", name, err))
	}

	outputSchema, err := jsonschema.For[Out](schemaOpts)
	if err != nil {
		panic(fmt.Sprintf("error inferring output schema for %s: %v", name, err))
	}

	tool := &mcp.Tool{
		Name:         name,
		Description:  desc,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}

	mcp.AddTool(s, tool, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, *Out, error) {
		out, err := fn(ctx, in)
		if err == nil {
			return nil, out, nil
		}
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		}, nil, nil
	})
}

func addObjectListTool[In, Out any](s *mcp.Server, name, desc string, fn func(context.Context, In) ([]Out, error)) {
	schemaOpts := &jsonschema.ForOptions{
		TypeSchemas: api.TypeSchemas(),
	}

	inputSchema, err := jsonschema.For[In](schemaOpts)
	if err != nil {
		panic(fmt.Sprintf("error inferring input schema for %s: %v", name, err))
	}

	outputSchema, err := jsonschema.For[listWrapper[Out]](schemaOpts)
	if err != nil {
		panic(fmt.Sprintf("error inferring output schema for %s: %v", name, err))
	}

	tool := &mcp.Tool{
		Name:         name,
		Description:  desc,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint: true,
		},
	}

	mcp.AddTool(s, tool, func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, *listWrapper[Out], error) {
		outList, err := fn(ctx, in)
		if err == nil {
			return nil, &listWrapper[Out]{Items: outList}, nil
		}
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		}, nil, nil
	})
}

type listWrapper[Out any] struct {
	Items []Out `json:"items"`
}
