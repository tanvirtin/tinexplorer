package main

import (
	"github.com/tanvirtin/tinexplorer/internal/argparser"
	"github.com/tanvirtin/tinexplorer/internal/db"
	"github.com/tanvirtin/tinexplorer/pkg/archiver"
	"log"
    "github.com/tanvirtin/tinexplorer/api/server"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	argparser := argparser.New()
	path := argparser.GetPath()

	if path != "" {
		gormDB, err := db.Create()
		checkErr(err)
		archiver := archiver.New(gormDB, 2500)
		err = archiver.Archive(path)
		checkErr(err)
	} else {
		_, err := db.Instance()
		checkErr(err)
	}

    server := server.New("localhost", 8080)
    err := server.Serve()
    checkErr(err)
}
