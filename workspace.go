package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	vault "github.com/sosedoff/ansible-vault-go"
)

type Entry struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Mode uint32 `json:"mode"`
	Time string `json:"time"`
	Data string `json:"data"`
}

type Workspace struct {
	user      string
	pass      string
	key       string
	localPath string
	storePath string
	entries   map[string]Entry
	logger    *log.Logger
	debug     bool
}

func workspaceForPath(path string, storePath string) *Workspace {
	return &Workspace{
		storePath: storePath,
		localPath: path,
		entries:   map[string]Entry{},
		logger:    log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (w *Workspace) writeEntriesFile() error {
	data, err := json.Marshal(w.entries)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.storePath, data, 0600)
}

func (w *Workspace) readEntriesFile() error {
	data, err := ioutil.ReadFile(w.storePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &w.entries)
}

func (w *Workspace) exists() bool {
	_, err := os.Stat(w.storePath)
	return err == nil
}

func (w *Workspace) init() error {
	if err := os.MkdirAll(filepath.Dir(w.storePath), 0750); err != nil {
		return err
	}
	return w.writeEntriesFile()
}

func (w *Workspace) list() ([]string, error) {
	items := []string{}
	for k := range w.entries {
		items = append(items, k)
	}
	return items, nil
}

// add adds a new file to the workspace
func (w *Workspace) add(path string) error {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	stat, err := os.Stat(fullpath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return err
	}

	encoded, err := vault.Encrypt(string(data), w.pass)
	if err != nil {
		return err
	}

	w.entries[fullpath] = Entry{
		Path: fullpath,
		Size: stat.Size(),
		Mode: uint32(stat.Mode()),
		Time: stat.ModTime().UTC().Format(time.RFC3339),
		Data: encoded,
	}

	return nil
}

// remove removes the file from the workspace
func (w *Workspace) remove(path string) error {
	if _, ok := w.entries[path]; !ok {
		return nil
	}
	delete(w.entries, path)

	return w.writeEntriesFile()
}

func (w *Workspace) fetch(path string) error {
	entry, ok := w.entries[path]
	if !ok {
		return errors.New("file is not tracked in workspace")
	}

	decoded, err := vault.Decrypt(entry.Data, w.pass)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, []byte(decoded), 0600); err != nil {
		return err
	}

	return os.Chmod(path, os.FileMode(entry.Mode))
}

func (w *Workspace) read(path string) (string, error) {
	entry, ok := w.entries[path]
	if !ok {
		return "", errors.New("file is not tracked in workspace")
	}

	decoded, err := vault.Decrypt(entry.Data, w.pass)
	if err != nil {
		return "", err
	}

	return decoded, nil
}

func (w *Workspace) destroy() error {
	return os.Remove(w.storePath)
}

func (w *Workspace) log(args ...interface{}) {
	if w.debug {
		w.logger.Println(args...)
	}
}
