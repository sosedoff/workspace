package main

import (
	"errors"
	"fmt"
)

// Runner runs commands
type Runner struct {
	before   func() error
	commands []Command
}

// Command has information about command and handler function
type Command struct {
	name        string
	description string
	alias       string
	handler     HandlerFunc
}

// HandlerFunc is a function to execute commands
type HandlerFunc func(*Workspace, []string) error

// NewRunner returns a new instance of a runner
func NewRunner() *Runner {
	return &Runner{
		commands: []Command{},
	}
}

func (r *Runner) findCommand(name string) *Command {
	for _, c := range r.commands {
		if c.name == name || c.alias == name {
			return &c
		}
	}
	return nil
}

func (r *Runner) commandExists(name string) bool {
	return r.findCommand(name) != nil
}

func (r *Runner) register(commands ...Command) {
	for _, cmd := range commands {
		if r.commandExists(cmd.name) {
			panic("command " + cmd.name + " already registered")
		}
		r.commands = append(r.commands, cmd)
	}
}

func (r *Runner) printUsage() {
	names := []string{}
	maxlen := 0

	for _, c := range r.commands {
		name := c.name
		if c.alias != "" {
			name = fmt.Sprintf("%s,%s", c.name, c.alias)
		}
		if len(name) > maxlen {
			maxlen = len(name)
		}
		names = append(names, name)
	}

	for i, name := range names {
		format := "%-" + fmt.Sprintf("%d", maxlen+2) + "s%s\n"
		fmt.Printf(format, name, r.commands[i].description)
	}
}

func (r *Runner) run(workspace *Workspace, args []string) error {
	if len(args) == 0 {
		return errors.New("command is required")
	}

	name := args[0]
	if name == "help" {
		r.printUsage()
		return nil
	}

	cmd := r.findCommand(name)
	if cmd == nil {
		return errors.New("invalid command")
	}

	if err := r.before(); err != nil {
		return err
	}

	return cmd.handler(workspace, args[1:])
}
