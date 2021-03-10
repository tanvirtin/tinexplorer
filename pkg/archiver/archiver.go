package archiver

import (
    "database/sql"
	"fmt"
    "path/filepath"
    "os"
    "path"
    "runtime"
    _ "github.com/mattn/go-sqlite3"
)

func getRootDir() string {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    return filepath.Dir(d)
}

func InitializeDatabase() error {
    rootDir := getRootDir();
    pathToDb := filepath.Join(rootDir, "../assets/tinexplore.db")

    os.Remove(pathToDb);

    db, err := sql.Open("sqlite3", pathToDb)

    if err != nil {
        return err
    }

    statement, err := db.Prepare(`
        CREATE TABLE IF NOT EXISTS file 
        (
            id INTEGER PRIMARY KEY, 
            name TEXT,
            path TEXT UNIQUE,
            extension TEXT NOT NULL,
            isDirectory INTEGER,
            parentDirectory TEXT,
            size INTEGER,
            createdDate TEXT NOT NULL,
            populatedDate TEXT NOT NULL
        )
    `)
    
    if err != nil {
        return err
    }

    if _, err := statement.Exec(); err != nil {
        return err
    }
    
    return nil
}

func Archive(path string) error {
    filepath.Walk(path, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if (fileInfo.IsDir()) {
            fmt.Printf("File -> %s -> %+v\n\n", path, fileInfo)
        } else {
            fmt.Printf("Directory -> %s -> %+v\n\n", path, fileInfo)
        }

        return nil
    })

    return nil
} 
