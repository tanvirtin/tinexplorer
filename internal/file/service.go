package file

import "gorm.io/gorm"

type Service struct {
    repository *Repository
}

func NewService(db *gorm.DB) *Service {
    return &Service{ repository: &Repository{ db: db }}
}

func (s *Service) Ls(path string) ([]File, error) {
    return s.repository.FindByParentDirectory(path)
}
