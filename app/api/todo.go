package api

import (
	"github.com/gin-gonic/gin"
	"mgo-gin/app/form"
	"mgo-gin/app/repository"
	"mgo-gin/db"
	err2 "mgo-gin/utils/err"
	"net/http"
)

func ApplyToDoAPI(app *gin.RouterGroup, resource *db.Resource) {
	toDoEntity := repository.NewToDoEntity(resource)
	toDoRoute := app.Group("/todos")

	toDoRoute.GET("", getAllToDo(toDoEntity))
	toDoRoute.GET("/:id", getToDoById(toDoEntity))
	toDoRoute.POST("", createToDo(toDoEntity))
	toDoRoute.PUT("/:id", updateToDo(toDoEntity))
}

// GetAllTodo godoc
// @Summary Get all Todotask
// @Description Get all Todotask
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} form.ToDoResp
// @Router /todos [get]
func getAllToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, err := toDoEntity.GetAll()
		response := map[string]interface{}{
			"todo": list,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// CreateToDo godoc
// @Summary Create Todotask
// @Description Create Todotask
// @Accept  json
// @Produce  json
// @Param todoForm body form.ToDoForm true "ToDoForm"
// @Success 200 {object} form.ToDoResp
// @Router /todos [post]
func createToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		todoReq := form.ToDoForm{}
		if err := ctx.Bind(&todoReq); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		todo, code, err := toDoEntity.CreateOne(todoReq)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetTodoById godoc
// @Summary Get Todotask by id
// @Description Get Todotask by id
// @Accept  json
// @Produce  json
// @Param id path int true "ToDotask ID"
// @Success 200 {object} form.ToDoResp
// @Router /todos/{id} [get]
func getToDoById(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		todo, code, err := toDoEntity.GetOneByID(id)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateToDo godoc
// @Summary Update Todotask
// @Description Update Todotask
// @Accept  json
// @Produce  json
// @Param todoForm body form.ToDoForm true "ToDoForm"
// @Param id path int true "ToDotask ID"
// @Success 200 {object} form.ToDoResp
// @Router /todos/{id} [put]
func updateToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		todoReq := form.ToDoForm{}
		if err := ctx.Bind(&todoReq); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		todo, code, err := toDoEntity.Update(id, todoReq)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
