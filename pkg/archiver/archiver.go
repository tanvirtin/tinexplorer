package archiver

import (
	"github.com/tanvirtin/tinexplorer/internal/file"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
    "strings"
)

type Archiver struct {
	fileRepository *file.Repository
}

func New(db *gorm.DB, batchSize int) *Archiver {
	fileRepository := file.NewRepository(db, batchSize)
	fileRepository.Sync()
	return &Archiver{fileRepository: fileRepository}
}

func createFile(id uint64, path string, rootPath string, fileInfo os.FileInfo) (*file.File, error) {
	parentDirectory, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	if filePath == parentDirectory {
		return nil, nil
	}

    filePath = strings.Replace(filePath, rootPath, ".", 1)
    parentDirectory = strings.Replace(parentDirectory, rootPath, ".", 1)

	return &file.File{
		ID:              id,
		Name:            fileInfo.Name(),
		Path:            filePath,
		Extension:       filepath.Ext(path),
		ParentDirectory: parentDirectory,
		Size:            fileInfo.Size(),
		IsDirectory:     fileInfo.IsDir(),
		CreatedDate:     fileInfo.ModTime().Unix(),
		PopulatedDate:   time.Now().Unix(),
	}, nil
}

func (a Archiver) Archive(rootPath string) error {
	var startTime time.Time = time.Now()
	var id uint64 = 0
	var count int = 0

	log.Printf("Archiving all files from path: %s", rootPath)

	if err := filepath.Walk(rootPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		id++
		file, err := createFile(id, path, rootPath, fileInfo)

		if file == nil {
			return nil
		}

		if err != nil {
			return err
		}

		if !a.fileRepository.Accumulate(*file) {
			if files, err := a.fileRepository.Flush(); err != nil {
				return err
			} else {
				count += len(files)
				log.Println("Records archived:", count)
			}
			a.fileRepository.Accumulate(*file)
		}

		return nil
	}); err != nil {
		return err
	} else {
		if files, err := a.fileRepository.Flush(); err != nil {
			return err
		} else {
			count += len(files)
			log.Println("Records archived:", count)
		}
	}

	log.Printf("Archiving %s took %s", rootPath, time.Since(startTime))

	return nil
}
