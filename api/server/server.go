package server

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/tanvirtin/tinexplorer/internal/file"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Server struct {
	db   *gorm.DB
	host string
	port int
}

func New(db *gorm.DB, host string, port int) *Server {
	return &Server{db: db, host: host, port: port}
}

func (s *Server) configureGraphqlResolvers() graphql.Fields {
	fileResolver := file.NewResolver(s.db)
	return graphql.Fields{
		"ls": fileResolver.Ls(),
	}
}

func (s *Server) Serve() error {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: s.configureGraphqlResolvers()}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		return err
	}

	http.Handle("/graphql", handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	}))

	log.Printf("Starting server on http://%s:%v", s.host, s.port)
	return http.ListenAndServe(fmt.Sprintf("%s:%v", s.host, s.port), nil)
}
