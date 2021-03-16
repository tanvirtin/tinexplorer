package server

import (
    "log"
    "fmt"
    "net/http"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type Server struct {
    host string
    port int
}

func New(host string, port int) *Server {
    return &Server{host: host, port: port}
}

func (s *Server) Serve() error {
	fields := graphql.Fields{
		"ls": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "directories", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

    if err != nil {
        return err
    }

	http.Handle("/graphql", handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: false,
        Playground: true,
	}))


    log.Printf("Server started on http://%s:%v", s.host, s.port)
	http.ListenAndServe(fmt.Sprintf("%s:%v", s.host, s.port), nil)
    
    return nil
}
