package deps

import (
	"os"

	"golang.org/x/mod/modfile"
)

type Golang struct {
	path    string
	modfile *modfile.File
}

func LoadGolang(path string) (*Golang, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modfile, err := modfile.Parse("go.mod", file, nil)
	if err != nil {
		return nil, err
	}

	return &Golang{
		path:    path,
		modfile: modfile,
	}, nil
}

func (g *Golang) Name() string {
	return "golang"
}

func (g *Golang) Version() string {
	return g.modfile.Go.Version
}

func (g *Golang) Install() error {
	return nil
}
