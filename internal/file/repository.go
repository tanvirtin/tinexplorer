package file

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
    batchSize int
    buffer []File
}

func NewRepository(db *gorm.DB, batchSize int) *Repository {
    return &Repository{ db: db, batchSize: batchSize, buffer: []File{} }
}

func (r *Repository) Sync() {
    r.db.AutoMigrate(&File{})
}

func (r *Repository) Push(model File) bool {
    if (len(r.buffer) == r.batchSize) {
        return false
    }
    r.buffer = append(r.buffer, model)
    return true
}

func (r *Repository) Flush() []File {
    flushedBuffer := r.buffer
    r.buffer = []File{}
    return flushedBuffer
}

func (r *Repository) BulkInsert(models []File) error {
    if result := r.db.CreateInBatches(models, len(models)); result.Error != nil {
        return result.Error;
    } 
    return nil 
}
