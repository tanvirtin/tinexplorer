package archiver

import (
	"os"
	"fmt"
	"log"
	"time"
	"gorm.io/gorm"
	"path/filepath"
	"github.com/tanvirtin/tinexplorer/internal/file"
)

type Archiver struct {
    fileRepository *file.Repository
    debug bool
}

func New(db *gorm.DB, batchSize int, debug bool) *Archiver {
    fileRepository := file.NewRepository(db, batchSize)
    fileRepository.Sync()
    return &Archiver { fileRepository: fileRepository, debug: debug }
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

func (a Archiver) log(str string) {
    if a.debug {
        log.Println(str)
    }
}

func (a Archiver) Archive(rootPath string) error {
    a.log(fmt.Sprintf("Archiving path: %s", rootPath))

    start := time.Now()
    var runningId uint64 = 0
    const concurrency int = 250
	totalInsertedRecords := 0
    channel := make(chan error, concurrency)

    err := filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        runningId++

        fileFile := createFileFile(runningId, path, fileInfo)

        if !a.fileRepository.Push(fileFile) {
            fileFiles := a.fileRepository.Flush()
            go func() {
                defer func() {
					<- channel
					totalInsertedRecords += len(fileFiles)
                    a.log(fmt.Sprintf("Records archived: %v", totalInsertedRecords))
				}()
                err := a.fileRepository.BulkInsert(fileFiles)
                channel <- err
            }()
            a.fileRepository.Push(fileFile)
        }

        return nil
    })

    if err != nil {
        return err
    }

    fileFiles := a.fileRepository.Flush()
    numRemainingFiles := len(fileFiles)

    if numRemainingFiles > 0 {
        go func() {
            defer func() {
                <- channel
                totalInsertedRecords += numRemainingFiles
                a.log(fmt.Sprintf("Records archived: %v", totalInsertedRecords))
            }()
            err := a.fileRepository.BulkInsert(fileFiles)
            channel <- err
        }()     
    }

    for i := 0; i < cap(channel); i++ {
        channel <- nil
    }

    a.log(fmt.Sprintf("Archive took: %v", time.Since(start)))

    return nil
}
