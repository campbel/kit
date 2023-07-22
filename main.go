package main

import (
	"os"
	"os/exec"

	"github.com/campbel/yoshi"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Taskfile struct {
	Version  string            `yaml:"version,omitempty"`
	Includes map[string]string `yaml:"includes,omitempty"`
	Tasks    map[string]any    `yaml:"tasks,omitempty"`
}

func main() {
	yoshi.New("kit").Run(func() error {
		data, err := os.ReadFile("taskfile.yml")
		if err != nil {
			return err
		}

		var taskfile Taskfile
		if err := yaml.Unmarshal(data, &taskfile); err != nil {
			return err
		}

		if err := os.MkdirAll(".kit", 0755); err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		includes := make(map[string]string)
		for k, v := range taskfile.Includes {
			path := cwd + "/.kit/" + k
			if _, err := os.Stat(path); err != nil {
				if err := Get(v, cwd+"/.kit/"+k, cwd, true); err != nil {
					return errors.Wrap(err, "failed to get kit "+k)
				}
				if taskfile.Includes == nil {
					taskfile.Includes = make(map[string]string)
				}
			}
			includes[k] = k
		}

		taskfile.Includes = includes
		out, err := yaml.Marshal(taskfile)
		if err != nil {
			return err
		}

		if err := os.WriteFile(".kit/taskfile.yml", out, 0644); err != nil {
			return err
		}

		cmd := exec.Command("task", append([]string{"--taskfile", ".kit/taskfile.yml"}, os.Args[1:]...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		return cmd.Run()
	})
}

func Get(src, dst, pwd string, dir bool) error {
	return (&getter.Client{
		Src:     src,
		Dst:     dst,
		Pwd:     pwd,
		Dir:     true,
		Options: nil,
	}).Get()
}
