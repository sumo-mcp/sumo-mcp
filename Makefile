sumo-mcp: bin
	go build -o bin/sumo-mcp cmd/main.go

bin:
	mkdir -p bin
