package main

import (
    "log"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/tanvirtin/tinexplorer/internal/args"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    args := args.New();
    path, err := args.GetPath()
    checkErr(err)

    err = archiver.Archive(path);
    checkErr(err)
}
