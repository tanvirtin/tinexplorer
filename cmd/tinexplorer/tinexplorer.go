package main

import (
    "os"
    "log"
    "path/filepath"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/akamensky/argparse"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func getPath() (string, error) {
	parser := argparse.NewParser("tinexplorer", "Filesystem GraphQL API ☁️")
	path := parser.String("p", "path", &argparse.Options { Required: true, Help: "Path used to archive your filesystem" })
	err := parser.Parse(os.Args)

	if err != nil {
        return "", err
	}

    absPath, err := filepath.Abs(*path)

    if err != nil {
        return "", err
    }

    return absPath, nil
}

func main() {
    path, err := getPath()
    checkErr(err)

    err = archiver.InitializeDatabase()
    checkErr(err)

    err = archiver.Archive(path);
    checkErr(err)
}
