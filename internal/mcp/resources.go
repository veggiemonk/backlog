package mcp

import (
	"context"
	_ "embed"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed prompt-cli.md
var PromptCLIInstructions string

//go:embed prompt-mcp.md
var PromptMCPInstructions string

// addResources adds all MCP resources to the server
func (s *Server) addResources() {
	geminiResource := &mcp.Resource{
		URI:         geminiInstructionsURI,
		Name:        "GEMINI.md",
		Description: "Instructions for how Gemini should use the backlog tool",
		MIMEType:    "text/markdown",
	}
	s.mcpServer.AddResource(geminiResource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      geminiInstructionsURI,
					MIMEType: "text/markdown",
					Text:     PromptCLIInstructions,
				},
			},
		}, nil
	})

	claudeResource := &mcp.Resource{
		URI:         claudeInstructionsURI,
		Name:        "CLAUDE.md",
		Description: "Instructions for how Claude should use the backlog tool",
		MIMEType:    "text/markdown",
	}
	s.mcpServer.AddResource(claudeResource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      claudeInstructionsURI,
					MIMEType: "text/markdown",
					Text:     PromptCLIInstructions,
				},
			},
		}, nil
	})

	agentResource := &mcp.Resource{
		URI:         agentInstructionsURI,
		Name:        "AGENTS.md",
		Description: "Instructions for how agents should behave in the backlog project",
		MIMEType:    "text/markdown",
	}
	s.mcpServer.AddResource(agentResource, func(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:      agentInstructionsURI,
					MIMEType: "text/markdown",
					Text:     PromptCLIInstructions,
				},
			},
		}, nil
	})
}
