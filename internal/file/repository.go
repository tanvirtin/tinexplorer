package file

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
    batchSize int
    buffer []Model
}

func NewRepository(db *gorm.DB, batchSize int) *Repository {
    db.AutoMigrate(&Model{})
    return &Repository{ db: db, batchSize: batchSize, buffer: []Model{} }
}

func (r *Repository) Push(model Model) bool {
    if (len(r.buffer) == r.batchSize) {
        return false
    }
    r.buffer = append(r.buffer, model)
    return true
}

func (r *Repository) Flush() []Model {
    flushedBuffer := r.buffer
    r.buffer = []Model{}
    return flushedBuffer
}

func (r *Repository) BulkInsert(models []Model) {
    r.db.CreateInBatches(models, len(models))
}
