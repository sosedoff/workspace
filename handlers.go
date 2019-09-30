package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func requireWorkspace(next HandlerFunc) HandlerFunc {
	return func(w *Workspace, args []string) error {
		if !w.exists() {
			return errWorkspaceMissing
		}
		return next(w, args)
	}
}

func requireFile(next HandlerFunc) HandlerFunc {
	return func(w *Workspace, args []string) error {
		if len(args) == 0 {
			return errFileRequired
		}
		return next(w, args)
	}
}

func handleInit(w *Workspace, args []string) error {
	if w.exists() {
		return errWorkspaceExists
	}
	return w.init()
}

func handleList(w *Workspace, args []string) error {
	for _, entry := range w.entries {
		fmt.Printf(
			"file: %v, size: %v, mod: %v\n",
			entry.Path,
			entry.Size,
			entry.Time,
		)
	}
	return nil
}

func handleAdd(w *Workspace, args []string) error {
	file, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		if err := w.add(file); err != nil {
			return err
		}
		return w.writeEntriesFile()
	}

	items := []string{}
	filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		items = append(items, path)
		return nil
	})

	fmt.Println("about to add the following files:")
	for _, item := range items {
		fmt.Println("-", item)
	}

	if !requireConfirmation("add files to workspace", "yes") {
		return errAbort
	}

	for _, item := range items {
		if err := w.add(item); err != nil {
			return err
		}
	}

	return w.writeEntriesFile()
}

func handleFetch(w *Workspace, args []string) error {
	matches := []string{}
	search := ""

	if len(args) > 0 {
		search = args[0]
	}

	for _, entry := range w.entries {
		shouldFetch := strings.HasPrefix(entry.Path, w.localPath)
		if shouldFetch && search != "" {
			shouldFetch = strings.Contains(entry.Path, search)
		}
		if shouldFetch {
			matches = append(matches, entry.Path)
		}
	}

	if len(matches) == 0 {
		fmt.Println("did not find any files to fetch")
		return nil
	}

	fmt.Println("about to fetch these files:")
	for _, entry := range matches {
		fmt.Println("-", entry)
	}

	if !requireConfirmation("continue", "yes") {
		return errAbort
	}

	for _, entry := range matches {
		if err := w.fetch(entry); err != nil {
			return err
		}
	}

	return nil
}

func handleRemove(w *Workspace, args []string) error {
	return w.remove(args[0])
}

func handleShow(w *Workspace, args []string) error {
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
	fmt.Println("* store path:", w.storePath)
	fmt.Println("* files tracked:", len(w.entries))
	return nil
}

func handleBackup(w *Workspace, args []string) error {
	originalPath := w.storePath
	defer func() {
		w.storePath = originalPath
	}()

	w.storePath = fmt.Sprintf("%s.backup.%v", originalPath, time.Now().Unix())

	fmt.Println("backing up store to", w.storePath)
	return w.writeEntriesFile()
}
