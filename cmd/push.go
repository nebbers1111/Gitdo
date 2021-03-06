package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Hands all tasks since last push to the plugin create function",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != nil {
			pDanger("Could not load gitdo: %v\n", err)
			return
		}
		if err := Push(cmd, args); err != nil {
			pDanger("Failed to run push: %v\n", err)
			return
		}

		pNormal("Gitdo finished pushing\n")
	},
}

// Push reads in tasks that are staged to be added, gives them to the create plugin and notifies the user that they were
// uploaded. Then moves them in to committed tasks and saves the task file. If the plugin fails, then the tasks are left
// and should be retried next 'git push'
func Push(cmd *cobra.Command, args []string) error {
	tasks, err := getTasksFile()
	if err != nil {
		return err
	}

	if len(tasks.NewTasks) == 0 && len(tasks.DoneTasks) == 0 {
		pInfo("No new tasks or done tasks\n")
		return nil
	}

	for id, task := range tasks.NewTasks {
		_, err = RunPlugin(CREATE, task)
		if err != nil {
			pDanger("Failed to add task '%s': %v\n", task.String(), err)
			continue
		}
		pInfo("Task %s added to %s\n", id, app.Plugin)
		tasks.RemoveTask(id)
	}

	failedIds := []string{}
	for _, id := range tasks.DoneTasks {
		_, err = RunPlugin(DONE, id)
		if err != nil {
			pWarning("Failed to mark %s as done\n", id)
			failedIds = append(failedIds, id)
			continue
		}
		pInfo("Task %s marked as done\n", id)
	}
	tasks.DoneTasks = failedIds

	err = writeTasksFile(tasks)
	if err != nil {
		return fmt.Errorf("could not save updated tasks list: %v", err)
	}

	return nil
}
