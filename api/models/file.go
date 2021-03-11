package models

type File struct {
    name string
    path string
    extension string
    parentDirectory string
    size uint64
    isDirectory bool
    createdDate string
    populatedDate string
}
