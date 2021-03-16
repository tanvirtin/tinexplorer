package argparser

import (
	"flag"
	"path/filepath"
	"strings"
)

type Args struct {
	path *string
}

func New() *Args {
	args := Args{
		path: flag.String("path", "", "Path used to archive your filesystem"),
	}
	flag.Parse()
	return &args
}

func (a *Args) GetPath() (string, error) {
	path := *a.path
	if strings.Contains(path, "~") {
		path = strings.ReplaceAll(path, "~", "/home")
	}
	return filepath.Abs(path)
}
