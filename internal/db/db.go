package db

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
)

func getDbPath() string {
    _, b, _, _ := runtime.Caller(0)
    d := path.Join(path.Dir(b))
    rootDir := filepath.Dir(d)
    pathToDb := filepath.Join(rootDir, "../assets/tinexplorer.db")
    return pathToDb
}

func Exists() bool {
    if _, err := os.Stat(getDbPath()); os.IsNotExist(err) {
        return false
    }
    return true
}

func Destroy() error {
    pathToDb := getDbPath()
    return os.Remove(pathToDb);
}

func Instance() (*gorm.DB, error) {
    pathToDb := getDbPath()

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
        return db, nil
    }
}
