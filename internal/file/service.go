package file

import (
	"fmt"
	"gorm.io/gorm"
)

type Service struct {
	repository *Repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repository: &Repository{db: db}}
}

func (s *Service) Ls(path string) ([]File, error) {
	return s.repository.FindByParentDirectory(path)
}

func (s *Service) Find(path string) ([]File, error) {
	return s.repository.QueryByPath(fmt.Sprintf("%s*", path))
}

func (s *Service) Search(text string) ([]File, error) {
	return s.repository.SearchByName(fmt.Sprintf("%%%s%%", text))
}
