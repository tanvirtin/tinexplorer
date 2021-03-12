package archiver

import (
    "os"
    "log"
    "path"
    "time"
    "runtime"
    "gorm.io/gorm"
    "path/filepath"
    "gorm.io/gorm/logger"
    "gorm.io/driver/sqlite"
    "github.com/tanvirtin/tinexplorer/api/models"
)

func getDbPath() string {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    rootDir := filepath.Dir(d)
    pathToDb := filepath.Join(rootDir, "../assets/tinexplore.db")
    return pathToDb
}

func createDb() (*gorm.DB, error) {
    pathToDb := getDbPath()
    os.Remove(pathToDb);

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Silent,
			Colorful: false,
		},
	)

    if db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{ Logger: newLogger }); err != nil {
        return nil, err
    } else {
        db.AutoMigrate(&models.File{})
        return db, nil
    }
}

func createFileModel(id uint64, path string, fileInfo os.FileInfo) models.File {
    return models.File{
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

func Archive(rootPath string) error {
    var runningId uint64 = 0
    const concurrency int = 250
    const batchSize int = 10000
	totalInsertedRecords := 0
    files := []models.File{}
    channel := make(chan bool, concurrency)

    db, err := createDb()

    if err != nil {
        return nil
    }

    err = filepath.Walk(rootPath, func (path string, fileInfo os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        runningId++

        file := createFileModel(runningId, path, fileInfo)

        if len(files) > batchSize {
            channel <- true
            go func() {
                defer func() {
					<-channel
					totalInsertedRecords += batchSize
                    log.Println("Records archived:", totalInsertedRecords)
				}()
                db.CreateInBatches(files, batchSize)
            }()
            files = []models.File{}
        }

        files = append(files, file);

        return nil
    })

    if err != nil {
        return err
    }

    channel <- true
    go func() {
		defer func() {
			<-channel
			totalInsertedRecords += len(files)
		}()
        db.CreateInBatches(files, batchSize)
    }()

    for i := 0; i < cap(channel); i++ {
        channel <- true
    }

    log.Println("Records archived:", totalInsertedRecords)

    return nil
} 
