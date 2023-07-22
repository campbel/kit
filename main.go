package main

import (
	"os"
	"os/exec"
	"reflect"

	"github.com/campbel/yoshi"
	"github.com/hashicorp/go-getter"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Taskfile struct {
	Version  string             `yaml:"version,omitempty"`
	Includes map[string]Include `yaml:"includes,omitempty"`
	Env      map[string]string  `yaml:"env,omitempty"`
	Dotenv   []string           `yaml:"dotenv,omitempty"`
	Tasks    yaml.MapSlice      `yaml:"tasks,omitempty"`
}

type Include struct {
	Taskfile string         `yaml:"taskfile,omitempty"`
	Dir      string         `yaml:"dir,omitempty"`
	Optional bool           `yaml:"optional,omitempty"`
	Vars     map[string]any `yaml:"vars,omitempty"`
	Aliases  []string       `yaml:"aliases,omitempty"`
}

func main() {
	yoshi.New("kit").Run(func() error {
		data, err := os.ReadFile("taskfile.yml")
		if err != nil {
			return err
		}

		var taskfile Taskfile
		if err := unmarshal(data, &taskfile); err != nil {
			return err
		}

		if err := os.MkdirAll(".kit", 0755); err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		includes := make(map[string]Include)
		for k, v := range taskfile.Includes {
			path := cwd + "/.kit/" + k
			if _, err := os.Stat(path); err != nil {
				if err := Get(v.Taskfile, cwd+"/.kit/"+k, cwd, true); err != nil {
					return errors.Wrap(err, "failed to get kit "+k)
				}
				if taskfile.Includes == nil {
					taskfile.Includes = make(map[string]Include)
				}
			}
			v.Taskfile = k
			includes[k] = v
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

func unmarshal(data []byte, v interface{}) error {
	var a any
	if err := yaml.Unmarshal(data, &a); err != nil {
		return errors.Wrap(err, "failed to unmarshal yaml")
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: func(fromType, toType reflect.Type, from any) (any, error) {
			// if fromType is string and toType is Include, map string to Include.Taskfile
			if fromType == reflect.TypeOf("") && toType == reflect.TypeOf(Include{}) {
				return Include{Taskfile: from.(string)}, nil
			}
			// if fromType is map and toType is yaml.MapSlice, map map to yaml.MapSlice
			if fromType == reflect.TypeOf(map[any]any{}) && toType == reflect.TypeOf(yaml.MapSlice{}) {
				var returnFrom yaml.MapSlice
				for k, v := range from.(map[any]any) {
					returnFrom = append(returnFrom, yaml.MapItem{Key: k, Value: v})
				}
				return returnFrom, nil
			}
			return from, nil
		},
		Result: &v,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create decoder")
	}

	return errors.Wrap(decoder.Decode(a), "failed to decode")
}
