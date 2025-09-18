package core

import (
	"slices"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestParseTaskWithDependencies(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog")

	// Create task A
	taskAParams := CreateTaskParams{
		Title:       "Task A",
		Description: "This is task A that will be a dependency for task B",
		AC:          []string{"Complete functionality A"},
	}

	taskA, err := store.Create(taskAParams)
	is.NoErr(err)
	is.Equal("Task A", taskA.Title)
	is.Equal("T01", taskA.ID.Name())

	// Create task B
	taskBParams := CreateTaskParams{
		Title:       "Task B",
		Description: "This is task B that depends on task A",
		AC:          []string{"Complete functionality B"},
	}

	taskB, err := store.Create(taskBParams)
	is.NoErr(err)
	is.Equal("Task B", taskB.Title)
	is.Equal("T02", taskB.ID.Name())

	// Create task C for testing multiple dependencies
	taskCParams := CreateTaskParams{
		Title:       "Task C",
		Description: "This is task C for additional dependency testing",
		AC:          []string{"Complete functionality C"},
	}

	taskC, err := store.Create(taskCParams)
	is.NoErr(err)
	is.Equal("Task C", taskC.Title)
	is.Equal("T03", taskC.ID.Name())

	t.Run("single dependency", func(t *testing.T) {
		// Add dependency from task B to task A using Update
		dependencies := []string{taskA.ID.Name()}
		editParams := EditTaskParams{
			ID:              taskB.ID.String(),
			NewDependencies: dependencies,
		}

		updatedTaskB, err := store.Update(taskB, editParams)
		is.NoErr(err)

		// Verify that task B now depends on task A
		is.Equal(len(updatedTaskB.Dependencies), 1)
		is.Equal(updatedTaskB.Dependencies[0], taskA.ID.Name())

		// Verify the dependency is persisted to file
		taskBFromFile, err := store.Get(taskB.ID.String())
		is.NoErr(err)
		is.Equal(len(taskBFromFile.Dependencies), 1)
		is.Equal(taskBFromFile.Dependencies[0], taskA.ID.Name())

		// Verify the file content contains the dependency in frontmatter
		filePath := ".backlog/T02-task_b.md"
		contentBytes, err := afero.ReadFile(fs, filePath)
		is.NoErr(err)
		content := string(contentBytes)

		// Check that dependencies are in the frontmatter
		is.True(strings.Contains(content, "dependencies:"))
		is.True(strings.Contains(content, "dependencies: T01"))
	})

	t.Run("multiple dependencies", func(t *testing.T) {
		// Add multiple dependencies to task C
		dependencies := []string{taskA.ID.Name(), taskB.ID.Name()}
		editParams := EditTaskParams{
			ID:              taskC.ID.String(),
			NewDependencies: dependencies,
		}

		updatedTaskC, err := store.Update(taskC, editParams)
		is.NoErr(err)

		// Verify that task C now depends on both task A and task B
		is.Equal(len(updatedTaskC.Dependencies), 2)
		is.True(slices.Contains(updatedTaskC.Dependencies.ToSlice(), taskA.ID.Name()))
		is.True(slices.Contains(updatedTaskC.Dependencies.ToSlice(), taskB.ID.Name()))

		// Verify the dependencies are persisted to file
		taskCFromFile, err := store.Get(taskC.ID.String())
		is.NoErr(err)
		is.Equal(len(taskCFromFile.Dependencies), 2)
		is.True(slices.Contains(taskCFromFile.Dependencies.ToSlice(), taskA.ID.Name()))
		is.True(slices.Contains(taskCFromFile.Dependencies.ToSlice(), taskB.ID.Name()))
	})
}

//
// func TestTaskSchema(t *testing.T) {
// 	is := is.New(t)
// 	customSchema, err := jsonschema.For[Task](nil)
// 	is.NoErr(err)
// 	// t.Log(customSchema)
// 	customSchema.Properties["id"].Type = "string"
// 	mcp.AddTool(server, &mcp.Tool{
// 		Name:        "customized greeting 2",
// 		InputSchema: customSchema,
// 	}, simpleGreeting)
// }
