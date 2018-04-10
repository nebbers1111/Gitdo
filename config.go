package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/urfave/cli"
	"strings"
)

type Config struct {
	// Author to attach to task in task manager.
	Author string `json:"author"`
	// Plugin to use at push time
	Plugin string `json:"plugin_name"`
	// The command to run for plugin files
	PluginInterpreter string `json:"plugin_interpreter"`

	// Example of plugin: "test" and plugin_interpreter: "python"
	// Will run 'python .git/gitdo/plugins/reserve_test'
}

// String returns a human readable format of the Config struct
func (c *Config) String() string {
	return fmt.Sprintf("Author: %s\nPlugin: %s\nInterpreter: %s", c.Author, c.Plugin, c.PluginInterpreter)
}

// Checks that the configuration has all the information needed
func (c *Config) IsSet() bool {
	if !c.pluginIsSet() {
		return false
	}
	if !c.authorIsSet() {
		return false
	}
	if !c.interpreterIsSet() {
		return false
	}
	return true
}

// pluginIsSet returns if the plugin in config is not empty
func (c *Config) pluginIsSet() bool {
	return strings.TrimSpace(c.Plugin) != ""
}

// authorIsSet returns if the author in config is not empty
func (c *Config) authorIsSet() bool {
	return strings.TrimSpace(c.Author) != ""
}

// interpreterIsSet returns if the plugin interpreter in config is not empty
func (c *Config) interpreterIsSet() bool {
	return strings.TrimSpace(c.PluginInterpreter) != ""
}

// LoadConfig opens a configuration file and reads it in to the Config struct
func LoadConfig(_ *cli.Context) error {
	bConfig, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		Warn("Could not find configuration file for gitdo. Have you ran \"gitdo init\"?")
		return err
	}

	err = json.Unmarshal(bConfig, config)
	if err != nil {
		return err
	}

	return nil
}

// getGitEmail runs the 'git config user.email' command to get the default email address of the user for the current dir
func getGitEmail() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	resp, err := cmd.Output()
	if err != nil {
		Warn("Please set your git email address for this repo. git config user.email example@email.com")
		return "", fmt.Errorf("Could not get user.email from git: %v", err)
	}
	return stripNewlineChar(resp), nil
}

// WriteConfig saves the current config to be loaded in after setting
func WriteConfig() error {
	bConf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFilePath, bConf, os.ModePerm)
	return err
}
