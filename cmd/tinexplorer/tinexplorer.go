package main

import (
    "log"
    "github.com/tanvirtin/tinexplorer/internal/db"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/tanvirtin/tinexplorer/internal/file"
    "github.com/tanvirtin/tinexplorer/internal/argparser"
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
        gormDB, err := db.Create()
        checkErr(err)
        archiver := archiver.New(gormDB, 2500)
        err = archiver.Archive(path)
        checkErr(err)
    } else {
        gormDB, err := db.Instance()
        checkErr(err)
        service := file.NewService(gormDB);
        files, err := service.Find("/home")
        checkErr(err)
        log.Println(files)
        files, err = service.Search("go.mod")
        checkErr(err)
        log.Println(files)
    }
}
