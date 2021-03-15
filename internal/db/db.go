package db

import (
    "os"
    "log"
    "path"
    "time"
    "errors"
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

func dbExists(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }
    return true
}

func openDb(path string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Silent,
			Colorful: false,
		},
	)

    if db, err := gorm.Open(sqlite.Open(path), &gorm.Config{ Logger: newLogger }); err != nil {
        return nil, err
    } else {
        return db, nil
    }
}

func Create() (*gorm.DB, error) {
    path := getDbPath()

    if dbExists(path) {
        err := Destroy()
        if err != nil {
            return nil, err
        }
    }

    return openDb(path)
}

func Instance() (*gorm.DB, error) {
    path := getDbPath()

    if !dbExists(path) {
        return nil, errors.New("Database does not exist") 
    }

    return openDb(path)
}

func Destroy() error {
    return os.Remove(getDbPath());
}
