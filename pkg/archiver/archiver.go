package archiver

import (
	"os"
	"time"
	"gorm.io/gorm"
	"path/filepath"
	"github.com/tanvirtin/tinexplorer/internal/file"
)

type Archiver struct {
    fileRepository *file.Repository
}

func New(db *gorm.DB, batchSize int) *Archiver {
    fileRepository := file.NewRepository(db, batchSize)
    fileRepository.Sync()
    return &Archiver { fileRepository: fileRepository }
}

func createFileFile(id uint64, path string, fileInfo os.FileInfo) file.File {
    return file.File{
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
}

func (a Archiver) Archive(rootPath string) error {
    var runningId uint64 = 0

    if err := filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        runningId++
        fileFile := createFileFile(runningId, path, fileInfo)

        if !a.fileRepository.Accumulate(fileFile) {
            if _, err := a.fileRepository.Flush(); err != nil {
                return err
            } else {
                a.fileRepository.Accumulate(fileFile)
            }
        }

        return nil
    }); err != nil {
        return err
    } else {
        if _, err := a.fileRepository.Flush(); err != nil {
            return err
        }
    }

    return nil
}
