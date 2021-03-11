package args

import (
    "os"
    "path/filepath"
    "github.com/akamensky/argparse"
)

type Args struct {
    parser *argparse.Parser
}

func New() *Args {
    return &Args{parser: argparse.NewParser("tinexplorer", "Filesystem GraphQL API ☁️")}
}

func (args *Args) GetPath() (string, error) {
	path := args.parser.String("p", "path", &argparse.Options { Required: true, Help: "Path used to archive your filesystem" })
	
	if err := args.parser.Parse(os.Args); err != nil {
        return "", err
	} 

    if absPath, err := filepath.Abs(*path); err != nil {
        return "", err
    } else {
        return absPath, nil
    }
}

