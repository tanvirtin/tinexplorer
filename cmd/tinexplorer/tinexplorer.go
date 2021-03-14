package main

import (
    "log"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/tanvirtin/tinexplorer/internal/argparser"
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

    argparser := argparser.New();
    path, err := argparser.GetPath()
    checkErr(err)

    if path != "" {
        archiver := archiver.New(db, 2700)
        err = archiver.Archive(path)
        checkErr(err)
    }
}
