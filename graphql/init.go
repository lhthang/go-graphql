package my_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RootFields= graphql.Fields{}

var ObjectID = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "id",
	Description:  "parse id",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case primitive.ObjectID:
			return value.Hex()
		default:
			return value
		}
		return nil
	},
	ParseValue:   nil,
	ParseLiteral: nil,
})
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
