package types

const TaskfileVersion3 = "3"

type Taskfile struct {
	Version string          `yaml:"version"`
	Tasks   map[string]Task `yaml:"tasks"`
}

type Task struct {
	Description string    `yaml:"desc,omitempty"`
	Internal    bool      `yaml:"internal,omitempty"`
	Status      []string  `yaml:"status,omitempty"`
	Commands    []Command `yaml:"cmds,omitempty"`
}

type Command struct {
	Command string         `yaml:"cmd,omitempty"`
	Task    string         `yaml:"task,omitempty"`
	Vars    map[string]any `yaml:"vars,omitempty"`
}
