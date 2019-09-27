package main

import (
	"errors"
)

var (
	errWorkspaceExists  = errors.New("workspace already exists")
	errWorkspaceMissing = errors.New("workspace is not configured")
	errFileRequired     = errors.New("file is required")
	errAbort            = errors.New("aborted")
)

func cmdInit() Command {
	return Command{
		name:        "init",
		description: "Init a new workspace",
		handler:     handleInit,
	}
}

func cmdList() Command {
	return Command{
		name:        "list",
		alias:       "ls",
		description: "List files in the current workspace",
		handler:     requireWorkspace(handleList),
	}
}

func cmdAdd() Command {
	return Command{
		name:        "add",
		description: "Add a new file or directory to the workspace",
		handler:     requireWorkspace(handleAdd),
	}
}

func cmdRemove() Command {
	return Command{
		name:        "remove",
		alias:       "rm",
		description: "Remove a file from the workspace",
		handler:     requireWorkspace(handleRemove),
	}
}

func cmdFetch() Command {
	return Command{
		name:        "fetch",
		description: "Fetch all files for the workspace",
		handler:     requireWorkspace(handleFetch),
	}
}

func cmdShow() Command {
	return Command{
		name:        "show",
		description: "Show file contents",
		handler:     requireWorkspace(handleShow),
	}
}

func cmdInfo() Command {
	return Command{
		name:        "info",
		description: "Show workspace info",
		handler:     requireWorkspace(handleInfo),
	}
}

func cmdDestroy() Command {
	return Command{
		name:        "destroy",
		description: "Destroy workspace",
		handler:     requireWorkspace(handleDestroy),
	}
}
