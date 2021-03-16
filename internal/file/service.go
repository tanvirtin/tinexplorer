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
	containsFirstPathDelimiter := string(path[0]) == "/"
	containsLastPathDelimiter := string(path[len(path)-1]) == "/"
	if !containsFirstPathDelimiter {
		path = fmt.Sprintf("/%s", path)
	}
	if containsLastPathDelimiter {
		path = path[0 : len(path)-1]
	}
	return s.repository.FindByParentDirectory(path)
}

func (s *Service) Find(path string) ([]File, error) {
	return s.repository.QueryByPath(fmt.Sprintf("%s*", path))
}

func (s *Service) Search(text string) ([]File, error) {
	return s.repository.SearchByName(fmt.Sprintf("%%%s%%", text))
}
