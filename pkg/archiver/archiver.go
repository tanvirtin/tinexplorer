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
    "time"
)

func Archive(rootPath string) error {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    rootDir := filepath.Dir(d)
    pathToDb := filepath.Join(rootDir, "../assets/tinexplore.db")

    os.Remove(pathToDb);

    db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{})
    
    if err != nil {
        return err
    }
        
    db.AutoMigrate(&models.File{})
 
    var id uint64 = 0
    filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        id += 1

        if (fileInfo.IsDir()) {
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
            db.Create(&file)
        } else {
             folder := models.File{
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
            db.Create(&folder)
        }

        return nil
    })

    return nil
} 
