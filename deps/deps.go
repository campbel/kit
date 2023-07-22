package deps

import "path/filepath"

type Dependency interface {
	Name() string
	Version() string
	Install() error
}

func DetermineDepsFromFiles(files []string) []Dependency {
	var deps []Dependency
	for _, path := range files {
		switch filepath.Base(path) {
		case "go.mod":
			dep, err := LoadGolang(path)
			if err != nil {
				continue
			}
			deps = append(deps, dep)
		}
	}
	return deps
}
