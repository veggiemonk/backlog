package mcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/imjasonh/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/logging"
)

const (
	geminiInstructionsURI = "mcp://backlog/GEMINI.md"
	claudeInstructionsURI = "mcp://backlog/CLAUDE.md"
	agentInstructionsURI  = "mcp://backlog/AGENTS.md"
)

// TaskStore interface matches the one expected by the MCP handlers
type TaskStore interface {
	Get(id string) (*core.Task, error)
	Create(params core.CreateTaskParams) (*core.Task, error)
	Update(task *core.Task, params core.EditTaskParams) (*core.Task, error)
	List(params core.ListTasksParams) ([]*core.Task, error)
	Search(query string, listParams core.ListTasksParams) ([]*core.Task, error)
	Path(t *core.Task) string
	Archive(id core.TaskID) (string, error)
}

// Server wraps the MCP server with backlog-specific functionality
type Server struct {
	mcpServer *mcp.Server
	handler   *handler
}

// handler contains the MCP tool implementations
type handler struct {
	store      TaskStore
	mu         *sync.Mutex
	autoCommit bool
}

// taskListResponse is used to wrap the list of tasks in a JSON object
// to conform to the MCP specification for structuredContent.
type taskListResponse struct {
	Tasks []*core.Task `json:"tasks"`
}

// NewServer creates a new MCP server configured for backlog
func NewServer(store TaskStore, autoCommit bool) *Server {
	ver := version.Get()
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "backlog MCP Server",
			Title:   "backlog",
			Version: ver.Version,
		}, &mcp.ServerOptions{
			Instructions: "Use this MCP server to manage your backlog tasks programmatically.",
			HasPrompts:   true,
			HasTools:     true,
			HasResources: true,
		},
	)

	h := &handler{
		store:      store,
		mu:         &sync.Mutex{},
		autoCommit: autoCommit,
	}

	server := &Server{
		mcpServer: mcpServer,
		handler:   h,
	}

	// Install all functionality
	server.addResources()
	server.addTools()
	server.addPrompts()

	return server
}

// RunHTTP starts the server with streamable HTTP transport
func (s *Server) RunHTTP(port int) error {
	addr := net.JoinHostPort("localhost", fmt.Sprintf("%d", port))
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server { return s.mcpServer }, nil)
	logging.Info("MCP server starting", "transport", "http", "address", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       10 * 1e9, // 10 seconds
		WriteTimeout:      10 * 1e9, // 10 seconds
		IdleTimeout:       60 * 1e9, // 60 seconds
		ReadHeaderTimeout: 5 * 1e9,  // 5 seconds
	}
	return server.ListenAndServe()
}

// RunStdio starts the server with stdio transport
func (s *Server) RunStdio(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}
