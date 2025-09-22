package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/veggiemonk/backlog/internal/core"
	"github.com/veggiemonk/backlog/internal/paths"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for issues in the backlog",
	Long:  `Scans the backlog for issues, such as duplicate or malformed task IDs.`, 
	Run: func(cmd *cobra.Command, args []string) {
		store, ok := cmd.Context().Value(ctxKeyStore).(TaskStore)
		if !ok {
			fmt.Println("Error: could not get task store from context")
			return
		}

		fmt.Println("Running the doctor...")

		fs := store.Fs()
		tasksDir := viper.GetString(configFolder)
		var err error
		tasksDir, err = paths.ResolveTasksDir(fs, tasksDir)
		if err != nil {
			fmt.Printf("Error resolving tasks directory: %v\n", err)
			return
		}

		files, err := afero.ReadDir(fs, tasksDir)
		if err != nil {
			fmt.Printf("Error reading tasks directory: %v\n", err)
			return
		}

		idMap := make(map[string]string)
		malformedFiles := []string{}
		duplicateIDs := make(map[string][]string)

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
				continue
			}

			// Check for malformed filenames
			if !strings.HasPrefix(file.Name(), "T") {
				malformedFiles = append(malformedFiles, file.Name())
				continue
			}

			// Read and parse the task file
			filePath := fmt.Sprintf("%s/%s", tasksDir, file.Name())
			content, err := afero.ReadFile(fs, filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			task, err := core.ParseTask(content)
			if err != nil {
				malformedFiles = append(malformedFiles, file.Name())
				continue
			}

			// Check for duplicate IDs
			if existingFile, ok := idMap[task.ID.String()]; ok {
				if _, ok := duplicateIDs[task.ID.String()]; !ok {
					duplicateIDs[task.ID.String()] = []string{existingFile}
				}
				duplicateIDs[task.ID.String()] = append(duplicateIDs[task.ID.String()], file.Name())
			} else {
				idMap[task.ID.String()] = file.Name()
			}
		}

		if len(malformedFiles) == 0 && len(duplicateIDs) == 0 {
			fmt.Println("No issues found. Your backlog is healthy!")
			return
		}

		if len(malformedFiles) > 0 {
			fmt.Println("\nMalformed files found:")
			for _, file := range malformedFiles {
				fmt.Printf("- %s\n", file)
			}
		}

		if len(duplicateIDs) > 0 {
			fmt.Println("\nDuplicate task IDs found:")
			for id, files := range duplicateIDs {
				fmt.Printf("- ID %s is duplicated in files: %s\n", id, strings.Join(files, ", "))
			}
		}
	},
}