# CLI Documentation Cleanup Analysis

**Goal**: Clean up CLI reference documentation after urfave/cli v3 migration to be easily understandable for users and AI agents.

## 1. Current Documentation Inventory

### Generated Documentation (Auto-generated via `make docs`)
Located in `docs/reference/`:
- `backlog.md` - Main CLI reference
- `backlog_archive.md`
- `backlog_create.md`
- `backlog_doctor.md`
- `backlog_edit.md`
- `backlog_instructions.md`
- `backlog_list.md`
- `backlog_mcp.md`
- `backlog_version.md`
- `backlog_view.md`

**Source**: Auto-generated from `internal/cmd/*.go` using `urfave/cli-docs` package in `internal/tools/docgen/docgen.go`

### Manual Documentation
- **README.md** - Main project documentation with CLI examples
- **internal/mcp/prompt-cli.md** - Comprehensive CLI instructions for AI agents (592 lines)
- **internal/mcp/prompt-mcp.md** - MCP server instructions for AI agents
- **docs/prompts/cli.md** - Processed copy of prompt-cli.md
- **docs/index.md** - Processed copy of README.md
- **docs/quick_start.md** - Extracted from README.md
- **docs/usage_examples.md** - Extracted from README.md
- **docs/ai_agent_integration.md** - Extracted from README.md

## 2. Key Issues Identified

### A. Duplication
1. **Section 9 in prompt-cli.md** ("Complete CLI Command Reference") duplicates auto-generated docs
   - Lines 342-433 contain manual CLI reference tables
   - This duplicates information already in `docs/reference/*.md`
   - **Risk**: Tables will become outdated as flags change

2. **README.md CLI examples** duplicate information in:
   - docs/usage_examples.md (extracted from README)
   - docs/quick_start.md (extracted from README)
   - prompt-cli.md workflow examples

### B. Inconsistent Flag Names (Post-Migration Issues)
Comparing actual CLI output vs documentation:

| Documentation Says | Actual CLI Flag | Location | Status |
|-------------------|----------------|----------|---------|
| `--assignee` | `--assigned, -a` | README.md line 410 | ❌ Wrong |
| `--unassign` | `--remove-assigned, -A` | prompt-cli.md line 142 | ❌ Wrong |
| Various flag descriptions | See actual `--help` | Multiple files | ⚠️ Needs verification |

### C. Terminology Issues
- Some references may still use Cobra-specific terminology
- Need to align all docs with urfave/cli v3 patterns
- Flag usage patterns need consistency (e.g., repeatable flags vs comma-separated)

### D. Structure Issues
1. **prompt-cli.md** is comprehensive (good!) but:
   - Section 9 should reference auto-generated docs instead of duplicating
   - Should maintain workflow examples but link to canonical reference

2. **README.md** should:
   - Link to detailed CLI reference instead of duplicating all flags
   - Keep high-level examples only

3. **Navigation**:
   - Users need clear path: README → CLI Reference → Detailed command docs
   - AI agents need: prompt-cli.md → CLI Reference → Command examples

## 3. Documentation Generation Flow

```
Source Code (internal/cmd/*.go)
    ↓
internal/tools/docgen/docgen.go
    ├─→ Generates docs/reference/*.md (via urfave/cli-docs)
    ├─→ Copies README.md → docs/index.md
    ├─→ Extracts sections from README.md → docs/{usage_examples,quick_start,ai_agent_integration}.md
    └─→ Copies internal/mcp/*.md → docs/prompts/*.md
```

**Key Insight**: The reference docs are **authoritative** and auto-generated. All other docs should reference them, not duplicate.

## 4. Actual CLI Structure (urfave/cli v3)

### Global Flags
```
--folder string      Directory for backlog tasks (default: ".backlog") [$BACKLOG_FOLDER]
--auto-commit        Auto-committing changes to git repository [$BACKLOG_AUTO_COMMIT]
--log-level string   Log level (debug, info, warn, error) (default: "info") [$BACKLOG_LOG_LEVEL]
--log-format string  Log format (json, text) (default: "text") [$BACKLOG_LOG_FORMAT]
--log-file string    Log file path (defaults to stderr) [$BACKLOG_LOG_FILE]
```

### Commands
- `archive` - Archive a task
- `create` - Create a new task
- `doctor` - Diagnose and fix task ID conflicts
- `edit` - Edit an existing task
- `instructions` - Instructions for agents to learn to use backlog
- `list` - List all tasks
- `mcp` - Start the MCP server
- `version` - Print the version information
- `view` - View a task by providing its ID

### Flag Pattern Differences from Documentation

**create command actual flags**:
- `--assigned, -a` (repeatable) - NOT `--assignee`
- `--labels, -l` (repeatable) - Can be used multiple times OR comma-separated
- `--ac` (repeatable) - Acceptance criteria

**edit command actual flags**:
- `--assigned, -a` (repeatable) - Add assigned names
- `--remove-assigned, -A` (repeatable) - NOT `--unassign`
- `--labels, -l` (repeatable) - Add labels
- `--remove-labels, -L` (repeatable) - Remove labels
- `--ac` (repeatable) - Add acceptance criterion
- `--check-ac` (repeatable) - Check AC by index
- `--uncheck-ac` (repeatable) - Uncheck AC by index
- `--remove-ac` (repeatable) - Remove AC by index

**list command actual flags**:
- `--status, -s` (repeatable) - Filter by status
- `--assigned, -a` (repeatable) - Filter by assigned names
- `--labels, -l` (repeatable) - Filter by labels

## 5. Proposed Cleanup Plan

### Phase 1: Fix Auto-Generated Documentation
1. ✅ Verify that `make docs` produces correct output
2. ✅ Ensure urfave/cli-docs is generating proper markdown
3. ✅ Check that all command descriptions are accurate

### Phase 2: Update prompt-cli.md (AI Agent Instructions)
1. **Keep** sections 1-8 (they provide valuable context for AI agents)
2. **Replace** section 9 "Complete CLI Command Reference" with:
   ```markdown
   ## 9. CLI Command Reference

   For the complete and authoritative CLI reference, see the [CLI Reference Documentation](../reference/backlog.md).

   Quick links to individual commands:
   - [backlog create](../reference/backlog_create.md)
   - [backlog edit](../reference/backlog_edit.md)
   - [backlog list](../reference/backlog_list.md)
   - [backlog view](../reference/backlog_view.md)
   - [backlog archive](../reference/backlog_archive.md)
   - [backlog doctor](../reference/backlog_doctor.md)

   ### Quick Reference Table

   | Task | Command Example |
   |------|----------------|
   | Create task | `backlog create "Title" --description "..." --ac "criterion"` |
   | Edit task | `backlog edit 42 --status "in-progress" --assigned "@you"` |
   | Check AC | `backlog edit 42 --check-ac 1` |
   | List tasks | `backlog list --status "todo" --assigned "alice"` |
   | View task | `backlog view 42` |
   | Archive task | `backlog archive 42` |
   ```

3. **Keep** sections 10-14 (pagination, workflows, common issues)
4. **Update** all flag names to match actual CLI:
   - `--assignee` → `--assigned`
   - `--unassign` → `--remove-assigned`
   - Verify all other flags

### Phase 3: Update README.md
1. **Simplify** CLI examples section - keep only high-level examples
2. **Add** prominent links to detailed docs:
   ```markdown
   For complete CLI reference, see [docs/reference/backlog.md](docs/reference/backlog.md)
   ```
3. **Remove** detailed flag tables (link to reference instead)
4. **Keep** conceptual examples and workflows

### Phase 4: Create CLI Navigation Guide
Create a new `docs/reference/README.md`:

```markdown
# CLI Reference

This directory contains the authoritative CLI reference documentation,
auto-generated from the source code.

## Quick Start
- New to backlog? Start with the [Quick Start Guide](../quick_start.md)
- Using backlog with AI agents? See [AI Agent Instructions](../prompts/cli.md)

## Command Reference
- [backlog](backlog.md) - Main command and global flags
- [backlog create](backlog_create.md) - Create new tasks
- [backlog edit](backlog_edit.md) - Edit existing tasks
- [backlog list](backlog_list.md) - List and filter tasks
- [backlog view](backlog_view.md) - View task details
- [backlog archive](backlog_archive.md) - Archive tasks
- [backlog doctor](backlog_doctor.md) - Diagnose and fix conflicts
- [backlog mcp](backlog_mcp.md) - MCP server
- [backlog instructions](backlog_instructions.md) - Get AI agent instructions
- [backlog version](backlog_version.md) - Version information

## Common Workflows
See [Usage Examples](../usage_examples.md) for common workflow patterns.
```

### Phase 5: Verification
1. Run `make docs` to regenerate all docs
2. Verify all examples with actual CLI:
   ```bash
   # Test each example from documentation
   ./bin/backlog create "Test" --assigned "test"  # Should work
   ./bin/backlog create "Test" --assignee "test"  # Should fail
   ```
3. Check all cross-links work correctly
4. Verify AI agent instructions still make sense

## 6. Files to Modify

### Must Change
- `internal/mcp/prompt-cli.md` - Remove section 9, fix flag names
- `README.md` - Simplify CLI sections, add links to reference
- Create: `docs/reference/README.md` - Navigation guide

### Verify and Update
- `docs/prompts/cli.md` - Will be regenerated from prompt-cli.md
- `docs/index.md` - Will be regenerated from README.md
- `docs/usage_examples.md` - Verify examples are correct
- `docs/quick_start.md` - Verify examples are correct
- `docs/ai_agent_integration.md` - Verify examples are correct

### Auto-Generated (Do Not Edit)
- `docs/reference/backlog*.md` - Generated by `make docs`

## 7. Success Criteria

✅ Auto-generated docs are the single source of truth for CLI reference
✅ No duplication between prompt-cli.md and reference docs
✅ All flag names match actual CLI implementation
✅ Clear navigation path: README → Reference → Detailed Commands
✅ AI agent instructions link to canonical reference
✅ All examples verified against actual CLI output
✅ Terminology consistent with urfave/cli v3
✅ Cross-links work correctly
✅ `make docs` regenerates everything correctly

## 8. Implementation Order

1. ✅ Create this analysis document
2. Fix `internal/mcp/prompt-cli.md` (section 9 + flag names)
3. Simplify README.md CLI sections
4. Create `docs/reference/README.md` navigation
5. Run `make docs` to regenerate
6. Verify all examples with actual CLI
7. Update any broken cross-links
8. Final review with actual CLI testing

## 9. Notes

- The documentation generator is at `internal/tools/docgen/docgen.go`
- It uses `github.com/urfave/cli-docs/v3` to generate markdown
- Command structure is defined in `internal/cmd/*.go`
- Generated docs use Jekyll front matter for GitHub Pages
- Keep AI agent instructions comprehensive - they're the primary interface for agents
