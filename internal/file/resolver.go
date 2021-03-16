package file

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

type Resolver struct {
	service  *Service
	fileType *graphql.Object
}

func NewResolver(db *gorm.DB) *Resolver {
	return &Resolver{
		service: NewService(db),
		fileType: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "File",
				Fields: graphql.Fields{
					"id":              &graphql.Field{Type: graphql.Int},
					"path":            &graphql.Field{Type: graphql.String},
					"name":            &graphql.Field{Type: graphql.String},
					"extension":       &graphql.Field{Type: graphql.Boolean},
					"parentDirectory": &graphql.Field{Type: graphql.String},
					"size":            &graphql.Field{Type: graphql.Int},
					"isDirectory":     &graphql.Field{Type: graphql.Boolean},
					"createdDate":     &graphql.Field{Type: graphql.Int},
					"populatedDate":   &graphql.Field{Type: graphql.Int},
				},
			},
		),
	}
}

func (r *Resolver) Ls() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(r.fileType),
		Description: "ls - List directory contents",
		Args:        graphql.FieldConfigArgument{"path": &graphql.ArgumentConfig{Type: graphql.String}},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			path, ok := p.Args["path"].(string)

			if !ok {
				return nil, nil
			}

			files, err := r.service.Ls(path)

			if err != nil {
				return nil, nil
			}

			return files, nil
		},
	}
}
