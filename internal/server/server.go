package server

import (
	"context"
	_ "embed"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/sumo-mcp/sumoapi-go"
)

//go:embed instructions.md
var instructions string

func New(version string, client sumoapi.Client) *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "sumo-mcp",
		Title:   "Sumo MCP Server",
		Version: version,
	}, &mcp.ServerOptions{
		Instructions: instructions,
		HasTools:     true,
	})

	addOpts := []mcp.AddToolOption{
		mcp.WithSchemaOptions(&jsonschema.ForOptions{
			TypeSchemas: sumoapi.TypeSchemas(),
		}),
	}

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "search_rikishi",
		Description: "Search for rikishi (sumo wrestlers).",
	}, wrap(client.SearchRikishi), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_rikishi",
		Description: "Get detailed profile information about a specific rikishi (sumo wrestler).",
	}, wrap(client.GetRikishi), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_rikishi_stats",
		Description: "Get overall performance stats about a specific rikishi (sumo wrestler).",
	}, wrap(client.GetRikishiStats), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_rikishi_matches",
		Description: "List matches for a specific rikishi (sumo wrestler).",
	}, wrap(client.ListRikishiMatches), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_rikishi_matches_against_opponent",
		Description: "List matches for a specific rikishi (sumo wrestler) against a specific opponent.",
	}, wrap(client.ListRikishiMatchesAgainstOpponent), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_basho",
		Description: "Get detailed information about a specific basho (sumo tournament).",
	}, wrap(client.GetBasho), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_banzuke",
		Description: "Get the banzuke (ranking list) for a specific basho (sumo tournament) division.",
	}, wrap(client.GetBanzuke), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_basho_with_torikumi",
		Description: "Get detailed information about a specific basho (sumo tournament) along with the torikumi (bout schedule) for a specific day and division.",
	}, wrap(client.GetBashoWithTorikumi), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_kimarite",
		Description: "List kimarite (winning techniques) along with their statistics.",
	}, wrap(client.ListKimarite), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_kimarite_matches",
		Description: "List matches won by a specific kimarite (winning technique).",
	}, wrap(client.ListKimariteMatches), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_measurement_changes",
		Description: "List measurement changes for rikishi (sumo wrestlers).",
	}, wrapList(client.ListMeasurementChanges), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_rank_changes",
		Description: "List rank changes for rikishi (sumo wrestlers).",
	}, wrapList(client.ListRankChanges), addOpts...)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_shikona_changes",
		Description: "List shikona (ring name) changes for rikishi (sumo wrestlers).",
	}, wrapList(client.ListShikonaChanges), addOpts...)

	return mcpServer
}

func wrap[In, Out any](fn func(context.Context, In) (Out, error)) func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, Out, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, Out, error) {
		out, err := fn(ctx, in)
		if err == nil {
			return nil, out, nil
		}
		var zero Out
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		}, zero, nil
	}
}

func wrapList[In, Out any](fn func(context.Context, In) ([]Out, error)) func(context.Context, *mcp.CallToolRequest, In) (*mcp.CallToolResult, *listWrapper[Out], error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, in In) (*mcp.CallToolResult, *listWrapper[Out], error) {
		outList, err := fn(ctx, in)
		if err == nil {
			return nil, &listWrapper[Out]{Items: outList}, nil
		}
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		}, nil, nil
	}
}

type listWrapper[Out any] struct {
	Items []Out `json:"items"`
}
