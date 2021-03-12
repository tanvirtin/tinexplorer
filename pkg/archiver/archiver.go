package archiver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tanvirtin/tinexplorer/internal/file"
	"gorm.io/gorm"
)

type Archiver struct {
    fileRepository *file.Repository
    debug bool
}

func New(db *gorm.DB, batchSize int, debug bool) *Archiver {
    return &Archiver { fileRepository: file.NewRepository(db, batchSize), debug: debug }
}

func createFileModel(id uint64, path string, fileInfo os.FileInfo) file.Model {
    return file.Model{
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
    channel := make(chan bool, concurrency)

    err := filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        runningId++

        fileModel := createFileModel(runningId, path, fileInfo)

        if !a.fileRepository.Push(fileModel) {
            fileModels := a.fileRepository.Flush()
            channel <- true
            go func() {
                defer func() {
					<-channel
					totalInsertedRecords += len(fileModels)
                    a.log(fmt.Sprintf("Records archived: %v", totalInsertedRecords))
				}()
                a.fileRepository.BulkInsert(fileModels)
            }()
            a.fileRepository.Push(fileModel)
        }

        return nil
    })

    if err != nil {
        return err
    }

    fileModels := a.fileRepository.Flush()
    numRemainingModels := len(fileModels)

    if numRemainingModels > 0 {
        channel <- true
        go func() {
            defer func() {
                <-channel
                totalInsertedRecords += numRemainingModels
            }()
            a.fileRepository.BulkInsert(fileModels)
        }()     
    }

    for i := 0; i < cap(channel); i++ {
        channel <- true
    }

    a.log(fmt.Sprintf("Records archived: %v", totalInsertedRecords))
    a.log(fmt.Sprintf("Archive took: %v", time.Since(start)))

    return nil
}
