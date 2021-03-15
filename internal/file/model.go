package file

import "gorm.io/gorm"

type File struct {
	gorm.Model
	ID              uint64 `gorm:"primaryKey"`
	Path            string `gorm:"unique"`
	Name            string
	Extension       string
	ParentDirectory string
	Size            int64
	IsDirectory     bool
	CreatedDate     int64
	PopulatedDate   int64
}
