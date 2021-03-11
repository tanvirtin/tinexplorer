package archiver

import (
	"fmt"
    "path/filepath"
    "os"
    "path"
    "runtime"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "github.com/tanvirtin/tinexplorer/api/models"
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

    if db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{}); err != nil {
        return err
    } else {
        db.AutoMigrate(&models.File{})
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
