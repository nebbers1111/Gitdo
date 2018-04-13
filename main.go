package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"path/filepath"

	"github.com/urfave/cli"
)

var (
	// Config needed for commit and push to use plugins and add author metadata
	// to task
	config = &Config{
		Author:            "",
		Plugin:            "",
		PluginInterpreter: "",
	}

	// Flags
	cachedFlag     bool
	verboseLogFlag bool

	// Current version
	version string

	// Gitdo working directory (holds plugins, secrets, tasks, etc.)
	gitdoDir = filepath.Join(".git", "gitdo")

	// File name for writing and reading staged tasks from (between commit
	// and push)
	stagedTasksFile = filepath.Join(gitdoDir, "tasks.json")
	configFilePath  = filepath.Join(gitdoDir, "config.json")
	pluginDirPath   = filepath.Join(gitdoDir, "plugins")
)

func main() {
	gitdo := AppBuilder()
	err := gitdo.Run(os.Args)
	if err != nil {
		Warnf("Gitdo failed to run: %v", err)
	}
}

// AppBuilder returns a urfave/cli app for directing commands and running setup
func AppBuilder() *cli.App {
	gitdo := cli.NewApp()
	gitdo.Name = "gitdo"
	gitdo.Usage = "track source code TODO comments - https://github.com/nebloc/Gitdo"
	gitdo.Version = "0.0.0-A5"
	if version != "" {
		gitdo.Version = fmt.Sprintf("App: %s, Build: %s_%s", version, runtime.GOOS, runtime.GOARCH)
	}
	gitdo.Before = ChangeToGitRoot
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print the app version",
	}
	gitdo.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose, v",
			Usage:       "sets logging to debug level",
			Destination: &verboseLogFlag,
		},
	}
	gitdo.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "prints the json of staged tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "config", Usage: "prints the current configuration"}},
			Action: List,
		},
		{
			Name:   "commit",
			Usage:  "gets git diff and stages any new tasks - normally ran from pre-commit hook",
			Action: Commit,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "cached, c",
					Usage:       "Diff is executed with --cached flag in commit mode",
					Destination: &cachedFlag,
				},
			},
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:   "init",
			Usage:  "sets the gitdo configuration for the current git repo",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "with-git", Usage: "Initialises a git repo first, then gitdo"}},
			Action: Init,
		},
		{
			Name:   "post-commit",
			Usage:  "adds the commit hash that has just been committed to tasks with empty hash fields",
			Action: PostCommit,
			After:  NotifyFinished,
		},
		{
			Name:   "push",
			Usage:  "starts the plugin to move staged tasks into your task manager - normally ran from pre-push hook",
			Action: Push,
			Before: LoadConfig,
			After:  NotifyFinished,
		},
		{
			Name:   "destroy",
			Usage:  "deletes all of the stored tasks",
			Flags:  []cli.Flag{cli.BoolFlag{Name: "yes", Usage: "confirms purge of task file"}},
			Before: ConfirmUser,
			Action: Destroy,
		},
	}
	return gitdo
}

// NotifyFinished prints that the process has finished and what command was ran
func NotifyFinished(ctx *cli.Context) error {
	log.Printf("Gitdo finished %s", ctx.Command.Name)
	return nil
}

// List pretty prints the tasks that are in file
func List(ctx *cli.Context) error {
	if ctx.Bool("config") {
		err := LoadConfig(ctx)
		if err != nil {
			return err
		}
		fmt.Println(config.String())
		return nil
	}
	tasks, _ := getTasksFile()

	fmt.Println(tasks.String())
	return nil
}

// stripNewLineChar takes a byte array (usually from an exec.Command run) and strips the newline characters, returning
// a string
func stripNewlineChar(orig []byte) string {
	var newStr string
	if strings.HasSuffix(string(orig), "\n") {
		newStr = string(orig)[:len(orig)-1]
	}
	if strings.HasSuffix(newStr, "\r") {
		newStr = newStr[:len(newStr)-1]
	}
	return newStr
}

// ChangeToGitRoot allows the running of Gitdo from subdirectories by moving the working dir to the top level according
// to git
func ChangeToGitRoot(_ *cli.Context) error {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	result, err := cmd.Output()
	if err != nil {
		// Not a git dir
		return errNotGitDir
	}
	err = os.Chdir(stripNewlineChar(result))
	return err
}
