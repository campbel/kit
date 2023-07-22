package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/campbel/yoshi"
	"gopkg.in/yaml.v2"
)

type Taskfile struct {
	Version  string            `yaml:"version,omitempty"`
	Kit      map[string]string `yaml:"kit,omitempty"`
	Includes map[string]string `yaml:"includes,omitempty"`
	Tasks    map[string]any    `yaml:"tasks,omitempty"`
}

type AnyMap map[string]any

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

		for k, v := range taskfile.Kit {
			resp, err := http.Get(v)
			if err != nil {
				return err
			}
			file, err := os.Create(".kit/" + k + ".yml")
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, resp.Body); err != nil {
				return err
			}
			if err := file.Close(); err != nil {
				return err
			}
			if err := resp.Body.Close(); err != nil {
				return err
			}
			if taskfile.Includes == nil {
				taskfile.Includes = make(map[string]string)
			}
			taskfile.Includes["kit:"+k] = k + ".yml"
		}

		taskfile.Kit = nil
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
