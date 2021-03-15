package argparser

import "flag"

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

func (a *Args) GetPath() string {
    return *a.path
}
