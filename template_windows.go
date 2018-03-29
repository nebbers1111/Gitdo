// +build windows

package main

import (
	"errors"
	"os/user"
	"path/filepath"
)

func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Could not determine user and therefore gitdo install directory")
	}
	gitdoPath := filepath.Join(usr.HomeDir, "AppData", "roaming", "Gitdo")
	return gitdoPath, nil
}
