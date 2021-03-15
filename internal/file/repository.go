package file

import "gorm.io/gorm"

type Repository struct {
	db        *gorm.DB
	batchSize int
	buffer    []File
}

func NewRepository(db *gorm.DB, batchSize int) *Repository {
	return &Repository{db: db, batchSize: batchSize, buffer: []File{}}
}

func (r *Repository) Sync() {
	r.db.AutoMigrate(&File{})
}

func (r *Repository) Accumulate(model File) bool {
	if len(r.buffer) == r.batchSize {
		return false
	}
	r.buffer = append(r.buffer, model)
	return true
}

func (r *Repository) Flush() ([]File, error) {
	buffer := r.buffer
	if result := r.db.CreateInBatches(r.buffer, len(r.buffer)); result.Error != nil {
		return buffer, result.Error
	} else {
		r.buffer = []File{}
		return buffer, nil
	}
}

func (r *Repository) FindByParentDirectory(directory string) ([]File, error) {
	files := []File{}
	if result := r.db.Where(&File{ParentDirectory: directory}).Find(&files); result.Error != nil {
		return nil, result.Error
	} else {
		return files, nil
	}
}

func (r *Repository) QueryByPath(path string) ([]File, error) {
	files := []File{}
	if result := r.db.Where("path GLOB ?", path).Find(&files); result.Error != nil {
		return nil, result.Error
	} else {
		return files, nil
	}
}

func (r *Repository) SearchByName(text string) ([]File, error) {
	files := []File{}
	if result := r.db.Where("name LIKE ?", text).Find(&files); result.Error != nil {
		return nil, result.Error
	} else {
		return files, nil
	}
}
