package main

import (
	"errors"

	"github.com/campbel/kit/commands"
	"github.com/campbel/kit/types"
	"github.com/campbel/yoshi"
)

type App struct {
	CD    func() error
	Edit  func() error
	Clone func() error
	Shell func(types.ShellOptions) error
}

func main() {
	yoshi.New("kit").Run(App{
		CD:    implementedInShell,
		Edit:  implementedInShell,
		Clone: implementedInShell,
		Shell: commands.Shell,
	})
}

func implementedInShell() error {
	return errors.New("only implemented in shell, eval `kit shell` to use")
}
