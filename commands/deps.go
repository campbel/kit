package commands

import (
	"fmt"
	"os"

	"github.com/campbel/kit/deps"
	"github.com/campbel/kit/types"
	"gopkg.in/yaml.v2"
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
	if opts.Task {
		taskfile := types.Taskfile{
			Version: types.TaskfileVersion3,
			Tasks:   make(map[string]types.Task),
		}
		for _, dep := range dependencies {
			taskfile.Tasks["setup:"+dep.Name()] = dep.Task()
		}
		data, err := yaml.Marshal(taskfile)
		if err != nil {
			return err
		}
		fmt.Print(string(data))
	} else {
		for _, dep := range dependencies {
			fmt.Println(dep.Name(), dep.Version())
		}
	}

	return nil
}
