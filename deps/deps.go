package deps

import (
	"path/filepath"

	"github.com/campbel/kit/types"
)

type Dependency interface {
	Name() string
	Version() string
	Task() types.Task
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
