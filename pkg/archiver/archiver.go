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

    db, err := gorm.Open(sqlite.Open(pathToDb), &gorm.Config{ Logger: newLogger })
    
    if err != nil {
        return nil, err
    }
        
    db.AutoMigrate(&models.File{})
    
    return db, nil
}

func Archive(rootPath string) error {
    start := time.Now()

    const concurrency int = 250
    const batchSize int = 10000

    db, err := createDb()

    if err != nil {
        return nil
    }

    var id uint64 = 0
	totalInsertedRecords := 0
    channel := make(chan bool, concurrency)

    files := []models.File{}
    
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

        if len(files) > batchSize {
            channel <- true
            go func() {
                defer func() {
					<-channel
					totalInsertedRecords += batchSize
					log.Println("Total records inserted ->", totalInsertedRecords)
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
			log.Println("Total records inserted ->", totalInsertedRecords)
		}()
        db.CreateInBatches(files, batchSize)
    }()

    for i := 0; i < cap(channel); i++ {
        channel <- true
    }

    elapsed := time.Since(start)

    log.Printf("Archive took %s", elapsed)

    return nil
} 
