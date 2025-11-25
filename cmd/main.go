package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/sumo-mcp/sumo-mcp/internal/api"
	"github.com/sumo-mcp/sumo-mcp/internal/server"
)

func main() {
	a := api.New(&http.Client{Timeout: 5 * time.Second})
	s := server.New("stdio", a)
	ctx := setupSignalHandler()
	if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func setupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
