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
    argparser := argparser.New();
    path, err := argparser.GetPath()
    checkErr(err)

    if path != "" {
        err := db.Destroy()
        checkErr(err)
    }

    db, err := db.Create()
    checkErr(err)
    archiver := archiver.New(db, 2700)

    if path != "" {
        err = archiver.Archive(path)
        checkErr(err)
    }
}
