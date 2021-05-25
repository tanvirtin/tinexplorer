package main

import (
    "github.com/tanvirtin/tinexplorer/api/server"
    "github.com/tanvirtin/tinexplorer/internal/argparser"
    "github.com/tanvirtin/tinexplorer/internal/db"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "gorm.io/gorm"
    "log"
)

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    var gormDB *gorm.DB
    var err error

    argparser := argparser.New()
    path, err := argparser.GetPath()
    checkErr(err)

    if path != "" {
        gormDB, err = db.Create()
        checkErr(err)
        archiver := archiver.New(gormDB, 2500)
        err = archiver.Archive(path)
        checkErr(err)
    } else {
        gormDB, err = db.Instance()
        checkErr(err)
    }

    server := server.New(gormDB, "localhost", 4000)
    err = server.Serve()
    checkErr(err)
}
