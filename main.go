package main

import (
	"os"
	"path/filepath"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		exitWithError(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		exitWithError(err)
	}
	storePath := filepath.Join(homeDir, "dropbox", "_workspace")

	workspace := workspaceForPath(dir, storePath)

	runner := NewRunner()
	runner.before = func() error {
		pass, err := getPassword()
		if err != nil {
			return err
		}
		workspace.pass = pass
		return nil
	}

	runner.register(
		cmdInit(),
		cmdFetch(),
		cmdAdd(),
		cmdRemove(),
		cmdShow(),
		cmdList(),
		cmdInfo(),
		cmdDestroy(),
	)

	if err := runner.run(workspace, os.Args[1:]); err != nil {
		exitWithError(err)
	}
}
