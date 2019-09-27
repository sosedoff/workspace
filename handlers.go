package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func requireWorkspace(next HandlerFunc) HandlerFunc {
	return func(w *Workspace, args []string) error {
		if !w.exists() {
			return errWorkspaceMissing
		}
		return next(w, args)
	}
}

func handleInit(w *Workspace, args []string) error {
	if w.exists() {
		return errWorkspaceExists
	}
	fmt.Println("setting up workspace for", w.localPath)
	return w.init()
}

func handleList(w *Workspace, args []string) error {
	files, err := w.list()
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f)
	}

	return nil
}

func handleAdd(w *Workspace, args []string) error {
	if len(args) == 0 {
		return errFileRequired
	}

	file := args[0]
	items := []string{}

	if file == "." {
		filepath.Walk(w.localPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			items = append(items, filepath.Base(path))
			return nil
		})
	} else {
		items = append(items, file)
	}

	for _, item := range items {
		fmt.Println("adding", item)
		if err := w.add(item); err != nil {
			return err
		}
		fmt.Println("added", item)
	}

	return nil
}

func handleFetch(w *Workspace, args []string) error {
	files, err := w.list()
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println("fetching", file)
		if err := w.fetch(file); err != nil {
			return err
		}
	}

	return nil
}

func handleRemove(w *Workspace, args []string) error {
	if len(args) == 0 {
		return errFileRequired
	}
	return w.remove(args[0])
}

func handleShow(w *Workspace, args []string) error {
	if len(args) == 0 {
		return errFileRequired
	}

	content, err := w.read(args[0])
	if err != nil {
		return err
	}

	fmt.Println(content)
	return nil
}

func handleDestroy(w *Workspace, args []string) error {
	if !requireConfirmation("destroy workspace", "yes") {
		return errAbort
	}

	if err := w.destroy(); err != nil {
		return err
	}

	fmt.Println("workspace has been destroyed")
	return nil
}

func handleInfo(w *Workspace, args []string) error {
	fmt.Println("workspace info:")
	fmt.Println("* local path:", w.localPath)
	fmt.Println("* storage path:", w.storePath)
	return nil
}
