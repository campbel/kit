package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

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
	Internal bool           `yaml:"internal,omitempty"`
	Vars     map[string]any `yaml:"vars,omitempty"`
	Aliases  []string       `yaml:"aliases,omitempty"`
}

var (
	IGNORE_CACHE = os.Getenv("KIT_IGNORE_CACHE") == "true"
)

func main() {
	// if no taskfile, call task anyways, let it handle the error
	taskfilePath, exists := getTaskFile(fileExists)
	if !exists {
		callTask(os.Args[1:])
		return
	}

	output, err := process(taskfilePath)
	if err != nil {
		panic(errors.Wrap(err, "kit failed to process taskfile"))
	}

	args := append(
		[]string{"--taskfile", output},
		filterArgs(os.Args)...,
	)

	callTask(args)
}

func callTask(args []string) error {
	cmd := exec.Command("task", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func process(taskfilePath string) (string, error) {
	data, err := os.ReadFile(taskfilePath)
	if err != nil {
		return "", err
	}

	var taskfile Taskfile
	if err := unmarshal(data, &taskfile); err != nil {
		return "", err
	}

	if err := os.MkdirAll(".kit", 0755); err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	includes := make(map[string]Include)
	for k, v := range taskfile.Includes {
		path := filepath.Join(cwd, ".kit", k)
		if !fileExists(path) || IGNORE_CACHE {
			if err := Get(v.Taskfile, path, cwd, true); err != nil {
				return "", errors.Wrap(err, "failed to get "+k)
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
		return "", err
	}

	outputPath := filepath.Join(cwd, ".kit", "taskfile.yml")
	if err := os.WriteFile(outputPath, out, 0644); err != nil {
		return "", err
	}

	return outputPath, nil
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

var (
	defaultTaskfiles = []string{
		"Taskfile.yml",
		"taskfile.yml",
		"Taskfile.yaml",
		"taskfile.yaml",
		"Taskfile.dist.yml",
		"taskfile.dist.yml",
		"Taskfile.dist.yaml",
		"taskfile.dist.yaml",
	}
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getTaskFile(fnFileExists func(path string) bool) (string, bool) {
	for i, arg := range os.Args[1:] {
		if (arg == "--taskfile" || arg == "-t") && len(os.Args) > i+2 {
			file := os.Args[i+2]
			return file, fileExists(file)
		}
		if strings.HasPrefix(arg, "--taskfile=") || strings.HasPrefix(arg, "-t=") {
			file := strings.Split(arg, "=")[1]
			return file, fileExists(file)
		}
	}
	for _, taskfile := range defaultTaskfiles {
		if fnFileExists(taskfile) {
			return taskfile, true
		}
	}
	return "", false
}

func filterArgs(args []string) []string {
	filteredArgs := []string{}
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if (arg == "--taskfile" || arg == "-t") && len(args) >= i+2 {
			i++
			continue
		}
		if strings.HasPrefix(arg, "--taskfile=") || strings.HasPrefix(arg, "-t=") {
			continue
		}
		filteredArgs = append(filteredArgs, arg)
	}
	return filteredArgs
}
