package my_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mgo-gin/app/form"
	"mgo-gin/app/model"
)

func populate() []form.ToDoResp {
	user := model.User{Id: primitive.NewObjectID(), Username: "Elliot Forbes"}
	todo:=model.ToDo{
		Id:       primitive.NewObjectID(),
		Name:     "Golang",
		Username: user.Username,
	}
	toDoResp := form.ToDoResp{
		ToDo: &todo,
		User: &user,
	}

	var tutorials []form.ToDoResp
	tutorials = append(tutorials, toDoResp)
	for _, todo := range tutorials {
		logrus.Printf("%+v", todo.ToDo)
		logrus.Printf("%+v", todo.User)
	}
	return tutorials
}

func Init() graphql.Schema  {
	toDoResps := populate()

	var toDoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "ToDo",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var typeTutorial = graphql.NewObject(graphql.ObjectConfig{
		Name: "ToDoResp",
		Fields: graphql.Fields{
			"todo": &graphql.Field{
				Type: toDoType,
			},
			"user": &graphql.Field{
				// here, we specify type as authorType
				// which we've already defined.
				// This is how we handle nested objects
				Type: userType,
			},
		},
	})

	fields := graphql.Fields{
		"todo": &graphql.Field{
			Name:        "",
			Type:        typeTutorial,
			Description: "Get todo by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				id, ok := p.Args["id"].(string)
				if ok {
					for _, tutorial := range toDoResps {
						if tutorial.ToDo.Id.Hex() == id {
							// return our tutorial
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(typeTutorial),
			Description: "Get Tutorial List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return toDoResps, nil
			},
		},
	}

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
