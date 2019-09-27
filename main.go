package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/keybase/go-keychain"
	"golang.org/x/crypto/ssh/terminal"
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
	items := []string{}

	if file == "." {
		filepath.Walk(workspace.localPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			items = append(items, filepath.Base(path))
			return nil
		})

		for _, item := range items {
			log.Println("---", item)
		}
	} else {
		items = append(items, file)
	}

	for _, item := range items {
		log.Println("adding", item, "to workspace")
		if err := workspace.add(item); err != nil {
			exitWithError(err)
		}
	}
}

func removeWorkspaceFile(workspace Workspace, file string) {
	if err := workspace.remove(file); err != nil {
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
		fmt.Println("fetching", file)
		if err := workspace.fetch(file); err != nil {
			exitWithError(err)
		}
	}
}

func showWorkspaceFile(workspace Workspace, file string) {
	content, err := workspace.read(file)
	if err != nil {
		exitWithError(err)
	}
	fmt.Println(content)
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
	case "ls", "list":
		requireWorkspace(workspace, listWorkspaceFiles)
	case "add":
		if len(args) < 2 {
			exitWithError("please provide a file name")
			return
		}
		requireWorkspace(workspace, func(workspace Workspace) {
			addWorkspaceFile(workspace, args[1])
		})
	case "rm", "remove":
		if len(args) < 2 {
			exitWithError("file required")
			return
		}
		requireWorkspace(workspace, func(workspace Workspace) {
			removeWorkspaceFile(workspace, args[1])
		})
	case "info":
		requireWorkspace(workspace, func(workspace Workspace) {
			fmt.Println("workspace info:")
			fmt.Println("local path:", workspace.localPath)
			fmt.Println("store path:", workspace.storePath)
		})
	case "show":
		if len(args) < 2 {
			exitWithError("please provide a file name")
			return
		}
		requireWorkspace(workspace, func(worker Workspace) {
			showWorkspaceFile(workspace, args[1])
		})
	case "fetch":
		requireWorkspace(workspace, fetchWorkspace)
	default:
		exitWithError("invalid command")
	}
}
