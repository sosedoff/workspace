package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/keybase/go-keychain"
)

func exitWithError(err interface{}) {
	fmt.Println("error:", err)
	os.Exit(1)
}

func requireWorkspace(workspace Workspace, next func(Workspace)) {
	if !workspace.exists() {
		exitWithError("workspace is not configured")
		return
	}
	next(workspace)
}

func listWorkspaceFiles(workspace Workspace) {
	files, err := workspace.list()
	if err != nil {
		exitWithError(err)
	}

	for _, f := range files {
		fmt.Println("->", f)
	}
}

func addWorkspaceFile(workspace Workspace, file string) {
	log.Println("adding", file, "to workspace")
	if err := workspace.add(file); err != nil {
		exitWithError(err)
	}
}

func initWorkspace(workspace Workspace) {
	if workspace.exists() {
		return
	}

	fmt.Println("setting up workspace for", workspace.localPath)
	if err := workspace.init(); err != nil {
		exitWithError(err)
	}
}

func configurePassword() (string, error) {
	service := "dotfiles"
	user := os.Getenv("USER")
	access := "com.dotfiles"

	pass, err := keychain.GetGenericPassword(service, user, "", access)
	if err != nil {
		return "", err
	}
	if len(pass) == 0 {
		fmt.Println("your password is not set")
		fmt.Printf("enter your new password: ")
		pass, err = terminal.ReadPassword(0)
		if err != nil {
			return "", err
		}
		item := keychain.NewGenericPassword(service, user, "", pass, access)
		if err := keychain.AddItem(item); err != nil {
			return "", err
		}
		fmt.Println("")
		return string(pass), nil
	}
	return string(pass), nil
}

func fetchWorkspace(workspace Workspace) {
	files, err := workspace.list()
	if err != nil {
		exitWithError(err)
		return
	}

	for _, file := range files {
		fmt.Println("adding", file)
		if err := workspace.fetch(file); err != nil {
			exitWithError(err)
		}
	}
}

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

	args := os.Args[1:]
	if len(args) == 0 {
		exitWithError("command is required")
	}

	pass, err := configurePassword()
	if err != nil {
		exitWithError(err)
	}

	workspace := workspaceForPath(dir, storePath)
	workspace.pass = pass

	switch args[0] {
	case "init":
		initWorkspace(workspace)
	case "list":
		requireWorkspace(workspace, listWorkspaceFiles)
	case "add":
		if len(args) < 2 {
			exitWithError("please provide a file name")
			return
		}
		requireWorkspace(workspace, func(workspace Workspace) {
			addWorkspaceFile(workspace, args[1])
		})
	case "fetch":
		requireWorkspace(workspace, fetchWorkspace)
	default:
		exitWithError("invalid command")
	}
}
