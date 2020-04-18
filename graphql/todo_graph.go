package my_graphql

import (
	"github.com/graphql-go/graphql"
	"mgo-gin/app/form"
	"mgo-gin/app/repository"
	"mgo-gin/utils/constant"
)


var toDoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ToDo",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: ObjectID,
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
			Type: ObjectID,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ToDoRespType = graphql.NewObject(graphql.ObjectConfig{
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

func InitToDoGraph() graphql.Fields {
	RootFields["getAllToDo"] = GetToDo()
	RootFields["getToDoById"] = GetToDoById()
	RootFields["createToDo"] = CreateToDo()
	return RootFields
}

func GetToDoById() *graphql.Field {
	return &graphql.Field{Name: "",
		Type:        ToDoRespType,
		Description: "Get todo by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
			id, ok := p.Args["id"].(string)
			if ok {
				toDo, _, err := repository.ToDoEntity.GetOneByID(id)
				if err != nil {
					return nil, err
				}
				return toDo, nil
			}
			return nil, nil
		},
	}
}

func GetToDo() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(ToDoRespType),
		Description: "Get all todo",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			list, _, err := repository.ToDoEntity.GetAll()
			if err != nil {
				return nil, err
			}
			return list, nil
		},
	}
}

func CreateToDo() *graphql.Field {
	return &graphql.Field{
		Type:        ToDoRespType,
		Description: "Create ToDo",
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"username": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			token:=p.Context.Value("token").(string)
			err :=RequireAuthorization(token,constant.ADMIN)
			if err!=nil{
				return nil,err
			}
			name, ok := p.Args["name"].(string)
			username, ok := p.Args["username"].(string)
			if ok {
				toDoForm := form.ToDoForm{
					Name:     name,
					Username: username,
				}
				toDo, _, err := repository.ToDoEntity.CreateOne(toDoForm)
				if err != nil {
					return nil, err
				}
				return toDo, nil
			}
			return nil, nil
		},
	}
}
