package deps

type Task struct {
	Name         string    `yaml:"name"`
	Description  string    `yaml:"desc"`
	Dependencies []string  `yaml:"deps"`
	Status       []string  `yaml:"status"`
	Commands     []Command `yaml:"cmds"`
}

type Command struct {
	Task string            `yaml:"task,omitempty"`
	Vars map[string]string `yaml:"vars,omitempty"`
}
