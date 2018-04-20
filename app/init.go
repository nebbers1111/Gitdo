package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nebloc/gitdo/app/utils"
	"github.com/urfave/cli"
)

// Init initialises the gitdo project by scaffolding the gitdo folder
func Init(ctx *cli.Context) error {
	if ctx.Bool("with-git") {
		if err := InitGit(); err != nil {
			return err
		}
	}

	utils.Highlightf("Making %s/gitdo", config.vc.NameOfDir())
	if err := os.MkdirAll(gitdoDir, os.ModePerm); err != nil {
		return err
	}

	if err := SetConfig(); err != nil {
		return err
	}

	if err := CreatePluginsDir(); err != nil {
		return err
	}

	utils.Highlight("Running plugin's setup...")
	if _, err := RunPlugin(SETUP, ""); err != nil {
		return err
	}

	if err := CreateHooks(); err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

// CreatePluginsDir creates a directory structure inside the Gitdo folder for Plugins to use as working space.
func CreatePluginsDir() error {
	path := filepath.Join(pluginDirPath, config.Plugin)
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

// SetConfig checks the config is not set and asks the user relevant questions to set it
func SetConfig() error {
	if config.IsSet() {
		return nil
	}

	if !config.authorIsSet() {
		author, err := AskAuthor()
		if err != nil {
			return err
		}
		config.Author = author
	}

	if !config.pluginIsSet() {
		plugin, err := AskPlugin()
		if err != nil {
			return err
		}
		config.Plugin = plugin
	}

	if !config.interpreterIsSet() {
		interp, err := GetInterp()
		if err != nil {
			utils.Warnf("No interp file in %s dir", config.Plugin)
			interp, err = AskInterpreter()
			if err != nil {
				return err
			}

		}
		config.PluginInterpreter = interp
	}

	err := WriteConfig()
	if err != nil {
		utils.Dangerf("Couldn't save config: %v", err)
		return err
	}
	return nil
}

// InitGit initialises a git repo before initialising gitdo
func InitGit() error {
	fmt.Println("Initializing git...")
	cmd := exec.Command("git", "init")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Println("Git initialized")
	return nil
}

// AskAuthor notifies user of email address used
func AskAuthor() (string, error) {
	email, err := config.vc.GetEmail()
	if err != nil {
		return "", err
	}
	utils.Highlightf("Using %s", email)
	return email, nil
}

// AskPlugin reads in plugins from the directory and gives the user a list of plugins, that have a "<name>_getid"
func AskPlugin() (string, error) {
	fmt.Println("Available plugins:")

	plugins, err := GetPlugins()
	if err != nil {
		return "", err
	}
	if len(plugins) < 1 {
		utils.Warn("No plugins found")
		return "", fmt.Errorf("no plugins")
	}
	for i, name := range plugins {
		fmt.Printf("%d: %s\n", i+1, name)
	}

	chosen := false
	pN := 0

	for !chosen {
		fmt.Printf("What plugin would you like to use (1-%d): ", len(plugins))
		var choice string
		_, err = fmt.Scan(&choice)
		if err != nil {
			return "", err
		}
		pN, err = strconv.Atoi(strings.TrimSpace(choice))
		if err != nil || pN > len(plugins) || pN < 1 {
			continue
		}
		chosen = true
	}
	plugin := plugins[pN-1]

	utils.Highlightf("Using %s", plugin)
	return plugin, nil
}

// AskInterpreter asks the user what command they want to use to run the plugin
func AskInterpreter() (string, error) {
	utils.Warn("Currently all plugins made as an example need python 3 set up in path. Redesign of plugin language choice and use coming soon.")
	var interp string
	for interp == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("What interpreter for this plugin (i.e. python3/node/python): ")
		var err error
		interp, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		interp = strings.TrimSpace(interp)
	}
	utils.Highlightf("Using %s", interp)
	return interp, nil
}

func GetInterp() (string, error) {
	homePath, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	contents, err := ioutil.ReadFile(filepath.Join(homePath, "plugins", config.Plugin, "interp"))
	if err != nil {
		return "", err
	}
	interp := utils.StripNewlineChar(contents)
	utils.Highlightf("Using %s - found in interp file", interp)
	return interp, err
}

// CreateHooks gets the users main Gitdo directory and copies the hooks from it to the correct version control hidden
// folder
func CreateHooks() error {
	homeDir, err := GetHomeDir()
	if err != nil {
		return err
	}
	utils.Highlight("Copying hooks...")
	return config.vc.SetHooks(homeDir)
}