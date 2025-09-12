---
id: "03"
title: Core Data Structures & Types
status: todo
labels:
    - core
    - types
    - data-structures
priority: high
created_at: 2025-09-12T05:33:45.241554Z
updated_at: 2025-09-12T06:30:10.428915Z
history:
    - timestamp: 2025-09-12T06:30:10.428914Z
      change: Implementation plan changed
---
## Description

Define the fundamental data structures that represent tasks, including proper YAML serialization, hierarchical ID system, and flexible type handling.

## Acceptance Criteria
<!-- AC:BEGIN -->

- [ ] #1 Task struct with complete YAML frontmatter support
- [ ] #2 TaskID type supporting hierarchical dot notation
- [ ] #3 MaybeStringArray for flexible string/array handling
- [ ] #4 Status and Priority enums with fuzzy matching validation

<!-- AC:END -->

## Implementation Plan

1. Define Task struct with all required fields and YAML tags\n2. Implement TaskID type with hierarchical parsing\n3. Create MaybeStringArray for flexible string/array handling\n4. Define Status and Priority enums with validation\n5. Add proper type methods and validation logic

## Implementation Notes


