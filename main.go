package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"Minimum Golang MCP Server",
		"1.0.0",
	)

	// add tool
	addTool := mcp.NewTool(
		"add",
		mcp.WithDescription("Add two numbers"),
		mcp.WithNumber("x",
			mcp.Required(),
		),
		mcp.WithNumber("y",
			mcp.Required(),
		),
	)
	s.AddTool(addTool, addToolHandler)

	// greeting resource
	greetingResource := mcp.NewResourceTemplate(
		"greeting://{name}",
		"getGreeting",
		mcp.WithTemplateDescription("Get a personalized greeting"),
		mcp.WithTemplateMIMEType("text/plain"),
	)
	s.AddResourceTemplate(greetingResource, greetingResourceHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func addToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	x := request.Params.Arguments["x"].(float64)
	y := request.Params.Arguments["y"].(float64)
	return mcp.NewToolResultText(fmt.Sprintf("%.2f", x+y)), nil
}

func greetingResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	name, err := extractNameFromURI(request.Params.URI)
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     fmt.Sprintf("Hello, %s!\n", name),
		},
	}, nil
}

// Extracts the name from a URI formatted as "greeting://{name}"
func extractNameFromURI(uri string) (string, error) {
	const prefix = "greeting://"
	if !strings.HasPrefix(uri, prefix) {
		return "", fmt.Errorf("invalid URI format: %s", uri)
	}
	name := strings.TrimPrefix(uri, prefix)
	if name == "" {
		return "", fmt.Errorf("name is empty in URI: %s", uri)
	}
	return name, nil
}
