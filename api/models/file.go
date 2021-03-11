package models

import (
    "gorm.io/gorm"
)

type File struct {
    gorm.Model
    ID uint64 `gorm:"primaryKey"`
    Name string 
    Path string `gorm:"unique"`
    Extension string
    ParentDirectory string
    Size int64
    IsDirectory bool
    CreatedDate int64
    PopulatedDate int64
}
