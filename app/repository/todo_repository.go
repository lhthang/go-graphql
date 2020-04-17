package repository

import (
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mgo-gin/app/form"
	"mgo-gin/app/model"
	"mgo-gin/db"
	"net/http"
)

var ToDoEntity IToDo

type toDoEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IToDo interface {
	GetAll() ([]form.ToDoResp, int, error)
	CreateOne(todo form.ToDoForm) (form.ToDoResp, int, error)
	GetOneByID(id string) (form.ToDoResp, int, error) // need return pointer
	Update(id string, todo form.ToDoForm) (form.ToDoResp, int, error)
}

//func NewToDoEntity
func NewToDoEntity(resource *db.Resource) IToDo {
	toDoRepo := resource.DB.Collection("todo")
	ToDoEntity = &toDoEntity{resource: resource, repo: toDoRepo}
	return ToDoEntity
}

func (entity *toDoEntity) GetAll() ([]form.ToDoResp, int, error) {
	toDoList := []form.ToDoResp{}
	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		return []form.ToDoResp{}, 400, err
	}

	for cursor.Next(ctx) {
		var todo model.ToDo
		err = cursor.Decode(&todo)
		if err != nil {
			logrus.Print(err)
			continue
		}
		user, _, err := UserEntity.GetOneByUsername(todo.Username)
		if err != nil {
			logrus.Print(err)
			continue
		}
		toDoList = append(toDoList, form.ToDoResp{
			ToDo: &todo,
			User: user,
		})
	}
	return toDoList, http.StatusOK, nil
}

func (entity *toDoEntity) CreateOne(todoForm form.ToDoForm) (form.ToDoResp, int, error) {
	user, _, err := UserEntity.GetOneByUsername(todoForm.Username)
	if user == nil || err != nil {
		return form.ToDoResp{}, getHTTPCode(err), err
	}
	todo := model.ToDo{
		Id:   primitive.NewObjectID(),
		Name: todoForm.Name,
	}
	ctx, cancel := initContext()
	defer cancel()
	_, err = entity.repo.InsertOne(ctx, todo)

	if err != nil {
		return form.ToDoResp{}, 400, err
	}
	todoResp := form.ToDoResp{
		ToDo: &todo,
		User: user,
	}
	return todoResp, http.StatusOK, nil
}

func (entity *toDoEntity) GetOneByID(id string) (form.ToDoResp, int, error) {
	var todo model.ToDo
	ctx, cancel := initContext()
	defer cancel()
	logrus.Print(id)
	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
	if err != nil {
		return form.ToDoResp{}, http.StatusNotFound, err
	}
	user, _, err := UserEntity.GetOneByUsername(todo.Username)
	if user == nil || err != nil {
		return form.ToDoResp{}, getHTTPCode(err), err
	}

	todoResp := form.ToDoResp{
		ToDo: &todo,
		User: user,
	}
	return todoResp, http.StatusOK, nil
}

func (entity *toDoEntity) Update(id string, todoForm form.ToDoForm) (form.ToDoResp, int, error) {
	var todo form.ToDoResp
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)

	todo, _, err := entity.GetOneByID(id)
	if err != nil {
		return form.ToDoResp{}, http.StatusNotFound, nil
	}

	user, _, err := UserEntity.GetOneByUsername(todoForm.Username)
	if user == nil || err != nil {
		return form.ToDoResp{}, getHTTPCode(err), err
	}

	err = copier.Copy(todo.ToDo, todoForm) // this is why we need return a pointer: to copy value
	if err != nil {
		logrus.Error(err)
		return form.ToDoResp{}, getHTTPCode(err), err
	}

	var newToDo model.ToDo
	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": todo}, opts).Decode(&newToDo)
	newResp := form.ToDoResp{
		ToDo: &newToDo,
		User: user,
	}
	return newResp, http.StatusOK, nil
}
