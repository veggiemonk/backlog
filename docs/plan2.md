# MCP Server Issues Fix Plan

## Executive Summary

After testing the backlog MCP functionality, three critical issues were identified that affect the reliability and usability of the MCP server. This document provides a comprehensive plan to address these issues systematically.

## Issues Identified

### 1. Response Size Limitation Error
**Issue**: `task_list` function returns responses exceeding 25,000 tokens, causing MCP client errors.
**Root Cause**: The function serializes all tasks at once without pagination or size limits.
**Impact**: High - Makes task listing unusable in repositories with many tasks.

### 2. Validation Errors in Task Edit Operations
**Issue**: Edit operations fail with validation errors related to field types (`acceptance_criteria`, `parent` field).
**Root Cause**: Inconsistent JSON schema validation and null/empty value handling.
**Impact**: High - Breaks core task editing functionality.

### 3. Inconsistent Field Type Handling
**Issue**: Inconsistency between null and empty array handling for `acceptance_criteria` field.
**Root Cause**: Go's JSON marshaling behavior for nil slices vs empty slices.
**Impact**: Medium - Causes unpredictable API responses.

## Detailed Fix Plan

### Phase 1: Response Size Management (Priority: Critical)

#### 1.1 Implement Pagination for task_list
**Files to modify**:
- `internal/mcp/tool_list.go`
- `internal/core/list.go` (if needed)

**Changes**:
1. Add pagination parameters to `ListTasksParams`:
   ```go
   type ListTasksParams struct {
       // existing fields...
       Limit  *int `json:"limit,omitempty"`
       Offset *int `json:"offset,omitempty"`
   }
   ```

2. Implement response size estimation before serialization
3. Add automatic pagination when response would exceed token limits
4. Provide pagination metadata in response

**Implementation Steps**:
1. Create a response size estimator function
2. Modify `handler.list()` to check response size before marshaling
3. If size exceeds limit, implement chunked response or return pagination instructions
4. Add tests for large task lists

#### 1.2 Add Response Size Monitoring
**Files to create/modify**:
- `internal/mcp/response_limiter.go` (new)
- `internal/mcp/tool_list.go`

**Changes**:
1. Create middleware to monitor response sizes
2. Add configurable size limits
3. Implement graceful degradation when limits are reached

### Phase 2: JSON Schema and Validation Fixes (Priority: Critical)

#### 2.1 Fix acceptance_criteria Field Handling
**Files to modify**:
- `internal/core/task.go`
- `internal/mcp/tool_edit.go`

**Root Cause Analysis**:
- Go marshals `nil` slices as `null` in JSON
- Go marshals empty slices (`[]AcceptanceCriterion{}`) as `[]` in JSON
- MCP validation expects consistent array type

**Solution**:
1. Ensure `AcceptanceCriteria` is always initialized as empty slice, never nil:
   ```go
   // In NewTask()
   AcceptanceCriteria: make([]AcceptanceCriterion, 0),
   ```

2. Add custom JSON marshaling for consistent behavior:
   ```go
   func (t *Task) MarshalJSON() ([]byte, error) {
       type Alias Task
       aux := (*Alias)(t)
       if aux.AcceptanceCriteria == nil {
           aux.AcceptanceCriteria = make([]AcceptanceCriterion, 0)
       }
       return json.Marshal(aux)
   }
   ```

#### 2.2 Fix parent Field Type Issues
**Files to modify**:
- `internal/core/task.go`
- `internal/core/id.go`

**Root Cause Analysis**:
- `TaskID` type may be marshaling inconsistently
- MCP schema expects string but may be receiving object

**Solution**:
1. Ensure `TaskID` consistently marshals to string:
   ```go
   func (t TaskID) MarshalJSON() ([]byte, error) {
       return json.Marshal(t.String())
   }
   ```

2. Handle empty TaskID consistently:
   ```go
   func (t TaskID) String() string {
       if len(t.seg) == 0 {
           return ""
       }
       // existing implementation
   }
   ```

#### 2.3 Add Comprehensive JSON Schema Validation
**Files to modify**:
- `internal/mcp/tool_edit.go`
- `internal/mcp/server.go`

**Changes**:
1. Enable commented-out JSON schema validation:
   ```go
   inputSchema, err := jsonschema.For[core.EditTaskParams](nil)
   if err != nil {
       return fmt.Errorf("jsonschema.For[core.EditTaskParams]: %v", err)
   }
   outputSchema, err := jsonschema.For[core.Task](nil)
   if err != nil {
       return fmt.Errorf("jsonschema.For[core.Task]: %v", err)
   }
   ```

2. Add proper error handling for validation failures
3. Implement schema validation middleware for all MCP tools

### Phase 3: Enhanced Error Handling (Priority: High)

#### 3.1 Improve MCP Error Messages
**Files to modify**:
- All `internal/mcp/tool_*.go` files
- `internal/mcp/server.go`

**Changes**:
1. Add structured error responses:
   ```go
   type MCPError struct {
       Code    int    `json:"code"`
       Message string `json:"message"`
       Details any    `json:"details,omitempty"`
   }
   ```

2. Wrap core errors with context for better debugging
3. Add error categorization (validation, not found, internal, etc.)

#### 3.2 Add Input Validation Middleware
**Files to create**:
- `internal/mcp/validation.go`

**Changes**:
1. Create validation middleware for all MCP tools
2. Validate task IDs before processing
3. Validate required fields and field formats
4. Return user-friendly validation error messages

### Phase 4: Performance and Scalability (Priority: Medium)

#### 4.1 Implement Smart Filtering
**Files to modify**:
- `internal/mcp/tool_list.go`
- `internal/core/list.go`

**Changes**:
1. Add response size estimation before querying
2. Implement intelligent default filters for large datasets
3. Add field selection to reduce response size
4. Cache commonly requested task lists

#### 4.2 Add Configurable Limits
**Files to modify**:
- `internal/mcp/server.go`
- `internal/config/` (new package)

**Changes**:
1. Add configurable response size limits
2. Add configurable pagination defaults
3. Add performance monitoring metrics

### Phase 5: Testing and Validation (Priority: High)

#### 5.1 Comprehensive MCP Testing
**Files to create/modify**:
- `internal/mcp/integration_test.go`
- `internal/mcp/large_dataset_test.go`

**Test Cases**:
1. Large task list scenarios (>1000 tasks)
2. Edge cases for all field types
3. Validation error scenarios
4. Response size limit testing
5. Pagination functionality

#### 5.2 Performance Testing
**Files to create**:
- `internal/mcp/performance_test.go`

**Test Scenarios**:
1. Response time with various task counts
2. Memory usage during large operations
3. Concurrent request handling

## Implementation Timeline

### Week 1: Critical Fixes
- [ ] Day 1-2: Fix acceptance_criteria and parent field JSON marshaling
- [ ] Day 3-4: Implement response size limiting for task_list
- [ ] Day 5: Add basic pagination support

### Week 2: Enhanced Error Handling
- [ ] Day 1-2: Improve error messages and validation
- [ ] Day 3-4: Add input validation middleware
- [ ] Day 5: Enable JSON schema validation

### Week 3: Testing and Polish
- [ ] Day 1-3: Comprehensive testing suite
- [ ] Day 4-5: Performance optimization and monitoring

## Success Criteria

1. **Response Size**: `task_list` operations complete successfully with any number of tasks
2. **Edit Operations**: All task edit operations work without validation errors
3. **Data Consistency**: All API responses have consistent field types
4. **Performance**: Response times remain under 2 seconds for datasets up to 10,000 tasks
5. **Error Handling**: All error cases return clear, actionable error messages

## Risk Mitigation

1. **Backward Compatibility**: Ensure all changes maintain API compatibility
2. **Data Migration**: Implement data validation for existing task files
3. **Rollback Plan**: Maintain ability to rollback to previous MCP implementation
4. **Testing**: Comprehensive test coverage before deployment

## Dependencies

- Go JSON schema validation library
- Existing core task management functionality
- MCP SDK compatibility

## Post-Implementation Monitoring

1. Monitor response sizes and pagination usage
2. Track error rates for different operations
3. Monitor performance metrics
4. Collect user feedback on error message clarity

---

*This plan addresses the immediate issues while establishing a foundation for long-term MCP server reliability and scalability.*