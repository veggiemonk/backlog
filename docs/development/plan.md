# Backlog Project Recreation Plan

This document outlines a comprehensive plan to recreate the backlog CLI task management tool from scratch, breaking the development into 10 major phases with detailed acceptance criteria.

## Project Overview

The backlog tool is a zero-configuration task manager written in Go that stores tasks as Markdown files in a Git repository. It features:

- **Hierarchical task structure** with dot notation IDs (T01 → T01.01 → T01.01.01)
- **Git-based storage** with automatic commits
- **AI-friendly MCP server** for agent collaboration
- **Comprehensive CLI interface** built with Cobra
- **Offline-first design** with complete portability

## Architecture Principles

- **File-based storage**: Tasks stored as Markdown files with YAML frontmatter in `.backlog/` directory
- **Hierarchical IDs**: Intelligent auto-generation with dot notation for parent-child relationships  
- **Interface-driven design**: TaskStore interface with afero filesystem abstraction for testability
- **Git integration**: Automatic commits for all task operations with meaningful messages
- **MCP integration**: Model Context Protocol server exposing task operations to AI agents
- **Type safety**: Custom types (TaskID, MaybeStringArray) with proper validation

## Development Phases

### Phase 1: Project Setup & Foundation (T01)

**Priority**: High  
**Labels**: foundation, setup

**Description**: Establish the basic project structure, build system, and development environment for the backlog CLI tool from scratch.

**Acceptance Criteria**:
- [ ] Go module initialized with proper naming (`github.com/veggiemonk/backlog`)
- [ ] Makefile configured with all build targets (build, test, lint, clean, docs, cover)
- [ ] Basic main.go entry point created with version handling
- [ ] Dependencies defined in go.mod with all required packages

**Key Dependencies**:
```go
github.com/spf13/cobra      // CLI framework
github.com/spf13/afero      // Filesystem abstraction
go.yaml.in/yaml/v4         // YAML processing
github.com/go-git/go-git/v5 // Git integration
github.com/modelcontextprotocol/go-sdk // MCP server
github.com/agnivade/levenshtein // Fuzzy matching
github.com/olekukonko/tablewriter // CLI tables
```

### Phase 2: Core Data Structures & Types (T03)

**Priority**: High  
**Labels**: core, types, data-structures

**Description**: Define the fundamental data structures that represent tasks, including proper YAML serialization, hierarchical ID system, and flexible type handling.

**Acceptance Criteria**:
- [ ] Task struct with complete YAML frontmatter support
- [ ] TaskID type supporting hierarchical dot notation parsing and validation
- [ ] MaybeStringArray for flexible string/array unmarshaling from YAML
- [ ] Status and Priority enums with fuzzy matching validation using Levenshtein distance

**Core Types**:
```go
type Task struct {
    ID           TaskID           `yaml:"id"`
    Title        string           `yaml:"title"`
    Status       Status           `yaml:"status"`
    Parent       *TaskID          `yaml:"parent,omitempty"`
    Assigned     MaybeStringArray `yaml:"assigned,omitempty"`
    Labels       MaybeStringArray `yaml:"labels,omitempty"`
    Dependencies MaybeStringArray `yaml:"dependencies,omitempty"`
    Priority     Priority         `yaml:"priority,omitempty"`
    CreatedAt    time.Time        `yaml:"created_at"`
    UpdatedAt    time.Time        `yaml:"updated_at"`
    History      []HistoryEntry   `yaml:"history,omitempty"`
}
```

### Phase 3: Task Storage System (T04)

**Priority**: High  
**Labels**: storage, filesystem, persistence

**Description**: Build the file-based storage system using afero filesystem abstraction for testability, with proper YAML frontmatter parsing and Markdown content handling.

**Acceptance Criteria**:
- [ ] TaskStore interface defined with all CRUD operations
- [ ] FileTaskStore implementation with afero filesystem for testability
- [ ] YAML frontmatter and Markdown parsing with proper separation
- [ ] Hierarchical task ID generation logic with intelligent auto-assignment

**File Structure**:
```
.backlog/
├── T01-project_setup.md          # Root task
├── T01.01-initialize_module.md   # Subtask
├── T01.01.01-create_gomod.md     # Sub-subtask
└── archived/                     # Archived tasks
```

**File Format**:
```markdown
---
id: "01.02.03"
title: "Implement OAuth integration"
status: "todo"
parent: "01.02"
assigned: ["alex", "jordan"]
labels: ["feature", "auth", "backend"]
priority: "high"
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-01-01T00:00:00Z
---

## Description
Task description here...

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Criterion one
- [x] #2 Criterion two (completed)
<!-- AC:END -->
```

### Phase 4: Core Business Logic (T05)

**Priority**: High  
**Labels**: business-logic, crud, search

**Description**: Implement the core task management operations including CRUD operations, search, filtering, and hierarchical task relationships.

**Acceptance Criteria**:
- [ ] Task creation with comprehensive validation and auto-ID assignment
- [ ] Task listing with filtering by status, parent, assignees, labels
- [ ] Task retrieval by ID with proper error handling
- [ ] Task update operations with field-specific updates (assignees, labels, AC)
- [ ] Content-based search functionality across task titles and descriptions
- [ ] Task archival system moving tasks to archived directory

**Key Operations**:
- `CreateTask(task Task) error`
- `ListTasks(filters FilterOptions) ([]Task, error)`
- `GetTask(id TaskID) (Task, error)`
- `UpdateTask(id TaskID, updates TaskUpdate) error`
- `SearchTasks(query string) ([]Task, error)`
- `ArchiveTask(id TaskID) error`

### Phase 5: CLI Commands & Interface (T06)

**Priority**: High  
**Labels**: cli, cobra, interface, commands

**Description**: Build the complete command-line interface using Cobra framework, providing all user-facing commands for task management.

**Acceptance Criteria**:
- [ ] Cobra CLI framework setup with proper command structure
- [ ] `create` command with full flag support (title, description, assignees, labels, priority, parent, AC)
- [ ] `list` command with filtering options (status, parent, assignee, labels)
- [ ] `view` command for detailed task display with formatted output
- [ ] `edit` command with field updates and add/remove operations
- [ ] `search` command with content-based search
- [ ] `archive` command for task archival
- [ ] `version` command with build information

**Command Examples**:
```bash
# Create comprehensive task
backlog create "Implement OAuth" \
  -d "Add Google OAuth integration" \
  -a "alice,bob" \
  -l "feature,auth" \
  --priority "high" \
  -p "T01" \
  --ac "Google OAuth works" \
  --ac "Token validation implemented"

# List with filters
backlog list --status "in-progress" --assigned "alice" --parent "T01"

# Edit with add/remove operations
backlog edit T01.02 \
  --status "in-progress" \
  --add-assignee "charlie" \
  --remove-label "draft" \
  --check-ac 1,3
```

### Phase 6: Git Integration (T07)

**Priority**: Medium  
**Labels**: git, integration, commits

**Description**: Implement automatic Git integration for task operations, ensuring all task changes are properly committed with meaningful messages.

**Acceptance Criteria**:
- [ ] Automatic Git commits for all task operations (create, update, archive)
- [ ] Proper commit message formatting following project conventions
- [ ] Git repository detection and validation with helpful error messages
- [ ] Error handling for Git operations with fallback behavior

**Commit Message Format**:
```
create: T01.02 Implement OAuth integration
update: T01.02 moved to in-progress, assigned alice
archive: T01.02 OAuth integration completed
```

### Phase 7: MCP Server Integration (T08)

**Priority**: High  
**Labels**: mcp, ai-integration, server

**Description**: Build the Model Context Protocol server to enable AI agent interaction with the task management system through standardized tools.

**Acceptance Criteria**:
- [ ] MCP server setup using the official go-sdk
- [ ] `task_create` tool with full parameter support
- [ ] `task_list` and `task_view` tools with proper formatting
- [ ] `task_edit` tool supporting all field updates
- [ ] `task_search` and `task_archive` tools
- [ ] HTTP transport support (`--http --port 8106`)
- [ ] STDIO transport support (default)
- [ ] Proper error handling and structured JSON responses

**MCP Tools**:
```go
// Available tools for AI agents
task_create    // Create new tasks with metadata
task_list      // List and filter existing tasks  
task_view      // Get detailed task information
task_edit      // Update task fields and properties
task_search    // Search tasks by content
task_archive   // Archive completed tasks
```

### Phase 8: Testing Infrastructure (T09)

**Priority**: High  
**Labels**: testing, unit-tests, coverage

**Description**: Establish comprehensive testing framework with unit tests, integration tests, and coverage reporting using afero filesystem mocking.

**Acceptance Criteria**:
- [ ] Testing framework setup with afero filesystem mocking for isolation
- [ ] Unit tests for all core task operations with edge cases
- [ ] Integration tests for CLI commands with real filesystem simulation
- [ ] MCP server functionality tests with tool validation
- [ ] Coverage reporting with HTML output and minimum thresholds

**Testing Patterns**:
```go
// Use afero.NewMemMapFs() for filesystem mocking
func TestCreateTask(t *testing.T) {
    fs := afero.NewMemMapFs()
    store := NewFileTaskStore(fs, ".backlog")
    
    task := Task{Title: "Test Task"}
    err := store.CreateTask(task)
    
    // Assertions...
}
```

### Phase 9: Documentation & Examples (T10)

**Priority**: Medium  
**Labels**: documentation, examples, guides

**Description**: Create comprehensive documentation including README, CLI docs, MCP integration guide, and development guidelines.

**Acceptance Criteria**:
- [ ] Comprehensive README with installation, usage examples, and feature overview
- [ ] Auto-generated CLI documentation using Cobra's doc generation
- [ ] MCP integration guide for AI agents with tool descriptions
- [ ] CLAUDE.md with development patterns, architecture notes, and coding standards
- [ ] Code examples and tutorials for common workflows

**Documentation Structure**:
```
docs/
├── cli/                 # Auto-generated CLI docs
├── plan.md             # This recreation plan
├── mcp-integration.md  # AI agent guide
└── architecture.md     # System design docs
```

### Phase 10: Polish & Release Preparation (T11)

**Priority**: Medium  
**Labels**: polish, release, optimization

**Description**: Final polishing including error handling, input validation, performance optimization, and comprehensive testing before release.

**Acceptance Criteria**:
- [ ] Consistent error handling patterns with context and actionable messages
- [ ] Input validation with helpful user feedback and suggestions
- [ ] Performance and memory optimization for large task sets
- [ ] Final integration testing across all components
- [ ] Release checklist completion with build verification

**Quality Standards**:
- All public functions have comprehensive error handling
- User inputs are validated with helpful error messages
- Memory usage is optimized for task sets of 1000+ items
- All edge cases are covered by tests
- Documentation is complete and up-to-date

## Development Workflow

### Prerequisites
- Go 1.25+ installed
- Git repository initialized
- Understanding of MCP (Model Context Protocol) concepts

### Build Commands
```bash
make build    # Build binary to bin/backlog
make test     # Run all tests with verbose output
make lint     # Run go vet linting
make cover    # Generate coverage report with HTML
make docs     # Generate CLI documentation
make clean    # Remove build artifacts
```

### Testing Strategy
1. **Unit Tests**: Test individual functions with afero filesystem mocking
2. **Integration Tests**: Test CLI commands end-to-end
3. **MCP Tests**: Validate MCP server tools and responses
4. **Coverage**: Maintain >80% test coverage across all packages

### Release Process
1. All phases completed with acceptance criteria met
2. Comprehensive testing including edge cases
3. Documentation review and updates
4. Performance benchmarking
5. Final integration testing
6. Version tagging and release notes

## Implementation Order

The phases should be implemented in order, as each builds upon the previous:

1. **T01** → Establishes foundation and build system
2. **T03** → Core types needed by all other components  
3. **T04** → Storage layer required for business logic
4. **T05** → Business logic powering CLI and MCP server
5. **T06** → User interface for manual task management
6. **T07** → Git integration for change tracking
7. **T08** → AI agent integration via MCP
8. **T09** → Quality assurance through testing
9. **T10** → User and developer documentation
10. **T11** → Final polish and optimization

Each phase can have internal parallelization, but dependencies between phases should be respected to avoid integration issues.

## Success Metrics

- [ ] CLI tool builds and runs without errors
- [ ] All task operations work correctly (CRUD + search)
- [ ] MCP server enables AI agent interaction
- [ ] Comprehensive test coverage (>80%)
- [ ] Complete documentation for users and developers
- [ ] Performance suitable for repositories with 1000+ tasks
- [ ] Zero-configuration setup for new users
