package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	vault "github.com/sosedoff/ansible-vault-go"
)

type Workspace struct {
	pass      string
	key       string
	localPath string
	storePath string
}

func workspaceForPath(path string, storePath string) Workspace {
	key := fmt.Sprintf("%x", sha1.Sum([]byte(path)))

	return Workspace{
		key:       key,
		storePath: filepath.Join(storePath, key),
		localPath: path,
	}
}

func (w Workspace) exists() bool {
	_, err := os.Stat(w.storePath)
	return err == nil
}

func (w Workspace) init() error {
	if err := os.MkdirAll(w.storePath, 0700); err != nil {
		return err
	}
	return nil
}

func (w Workspace) list() ([]string, error) {
	items := []string{}
	filepath.Walk(w.storePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		items = append(items, filepath.Base(path))
		return nil
	})
	return items, nil
}

func (w Workspace) add(file string) error {
	fullPath := filepath.Join(w.localPath, file)
	if _, err := os.Stat(fullPath); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(w.storePath, file)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0700); err != nil {
		return err
	}

	err = vault.EncryptFile(dstPath, string(data), w.pass)
	if err != nil {
		return err
	}
	return nil
}

func (w Workspace) fetch(file string) error {
	srcPath := filepath.Join(w.storePath, file)
	dstPath := filepath.Join(w.localPath, file)

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}

	raw, err := vault.Decrypt(string(data), w.pass)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dstPath, []byte(raw), 0644)
}

func (w Workspace) read(file string) (string, error) {
	srcPath := filepath.Join(w.storePath, file)

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return "", err
	}

	raw, err := vault.Decrypt(string(data), w.pass)
	if err != nil {
		return "", err
	}

	return raw, nil
}
