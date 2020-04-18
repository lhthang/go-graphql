package my_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var RootFields= graphql.Fields{}

func Init() graphql.Schema  {

	InitToDoGraph()
	fields := RootFields

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logrus.Print(err)
	}
	return schema
}
