package main

import (
	"fmt"
	"os"

	"github.com/keybase/go-keychain"
	"golang.org/x/crypto/ssh/terminal"
)

func getPassword() (string, error) {
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
