package main

import (
    "log"
    "github.com/tanvirtin/tinexplorer/pkg/archiver"
    "github.com/tanvirtin/tinexplorer/internal/argparser"
    "github.com/tanvirtin/tinexplorer/internal/file"
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
        if db.Exists() {
            err = db.Destroy()
            checkErr(err)
        }
        gormDB, err := db.Instance()
        checkErr(err)
        archiver := archiver.New(gormDB, 2500)
        err = archiver.Archive(path)
        checkErr(err)
    } 

    if db.Exists() {
        gormDB, err := db.Instance()
        checkErr(err)
        service := file.NewService(gormDB);
        if files, err := service.Find("/home/tanvirtin/workspace/tinexplorer/go.mod"); err != nil {
            log.Fatal(err)
        } else {
            log.Println(files)
        }
    } else {
        log.Fatal("Database does not exist")
    }
}
