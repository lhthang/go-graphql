package form

import "mgo-gin/app/model"

type ToDoForm struct {
	Name   string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type ToDoResp struct {
	*model.ToDo
	*model.User
}
