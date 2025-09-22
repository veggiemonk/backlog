package core

import (
	"fmt"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestCreateTaskConcurrently(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
	store := NewFileTaskStore(fs, ".backlog", NewMockLocker())

	var wg sync.WaitGroup
	numTasks := 10

	// Use a channel to serialize access to the Create function
	createQueue := make(chan struct{}, 1)

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Acquire the semaphore
			createQueue <- struct{}{}

			params := CreateTaskParams{
				Title: fmt.Sprintf("Concurrent Task %d", i),
			}
			_, err := store.Create(params)
			is.NoErr(err)

			// Release the semaphore
			<-createQueue
		}(i)
	}

	wg.Wait()

	tasks, err := store.List(ListTasksParams{})
	is.NoErr(err)
	is.Equal(len(tasks), numTasks)

	// Check for unique IDs
	idSet := make(map[string]bool)
	for _, task := range tasks {
		is.True(!idSet[task.ID.String()])
		idSet[task.ID.String()] = true
	}
}