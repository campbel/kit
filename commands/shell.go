package commands

import (
	"fmt"
	"strings"

	"github.com/campbel/kit/shell"
	"github.com/campbel/kit/types"
)

func Shell(opts types.ShellOptions) error {
	shell, ok := shellMap[opts.Shell]
	if !ok {
		return fmt.Errorf("unknown shell %s, must be one of %s", opts.Shell, strings.Join(shells, ","))
	}
	fmt.Println(shell)
	return nil
}

var shells = (func() []string {
	shells := make([]string, 0, len(shellMap))
	for k := range shellMap {
		shells = append(shells, k)
	}
	return shells
})()

var shellMap = map[string]string{
	"zsh": shell.ZSH,
}
