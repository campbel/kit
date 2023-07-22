package commands

import (
	"fmt"
	"os"

	"github.com/campbel/kit/deps"
	"github.com/campbel/kit/types"
)

func Deps(opts types.DepsOptions) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, dir+"/"+entry.Name())
	}

	dependencies := deps.DetermineDepsFromFiles(files)
	for _, dep := range dependencies {
		fmt.Println(dep.Name(), dep.Version())
	}

	return nil
}
