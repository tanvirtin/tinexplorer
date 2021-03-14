package main

import (
    "log"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/tanvirtin/tinexplorer/internal/args"
    "github.com/tanvirtin/tinexplorer/internal/db"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    db, err := db.Create()
    checkErr(err)

    args := args.New();
    path, err := args.GetPath()
    checkErr(err)

    archiver := archiver.New(db, 2700)
    err = archiver.Archive(path)
    checkErr(err)
}
