package archiver

import (
    "path/filepath"
    "os"
    "path"
    "runtime"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "github.com/tanvirtin/tinexplorer/api/models"
    "time"
    "log"
)

func getDbPath() string {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    rootDir := filepath.Dir(d)
    pathToDb := filepath.Join(rootDir, "../assets/tinexplore.db")
    return pathToDb
}

func Archive(rootPath string) error {
    start := time.Now()

    pathToDb := getDbPath()
    os.Remove(pathToDb);

    db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{})
    
    if err != nil {
        return err
    }
        
    db.AutoMigrate(&models.File{})
 
    var id uint64 = 0
    channel := make(chan bool, 100)
    
    err = filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        id++

        file := models.File{
            ID: id,
            Name: fileInfo.Name(),
            Path: path,
            Extension: filepath.Ext(path),
            ParentDirectory: filepath.Dir(path),
            Size: fileInfo.Size(),
            IsDirectory: fileInfo.IsDir(),
            CreatedDate: fileInfo.ModTime().Unix(),
            PopulatedDate: time.Now().Unix(),
        }

        channel <- true
        go func() {
            defer func() { <-channel }()
            db.Create(&file)
            log.Println("Created file ->", path);
        }()

        return nil
    })

    if err != nil {
        return err
    }

    for i := 0; i < cap(channel); i++ {
        channel <- true
    }

    elapsed := time.Since(start)

    log.Printf("Archive took %s", elapsed)

    return nil
} 
