package main

import (
	"os"
	"path/filepath"
)

func getStorePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "dropbox/workspace/config.json"), nil
}

func main() {
	storePath, err := getStorePath()
	if err != nil {
		exitWithError(err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		exitWithError(err)
		return
	}

	workspace := workspaceForPath(dir, storePath)

	runner := NewRunner()
	runner.before = func() error {
		pass, err := getPassword()
		if err != nil {
			return err
		}
		workspace.pass = pass

		if workspace.exists() {
			if err := workspace.readEntriesFile(); err != nil {
				return err
			}
		}

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
		cmdBackup(),
	)

	if err := runner.run(workspace, os.Args[1:]); err != nil {
		exitWithError(err)
	}
}
