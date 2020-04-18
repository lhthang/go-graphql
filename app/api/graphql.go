package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/graphiql"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
)

type reqBody struct {
	Query string `json:"query"`
}

func ApplyGraphAPI(app *gin.RouterGroup, schema *graphql.Schema) {

	route := app.Group("/graphql")
	route.Any("", handleGraph(*schema))

	graphiqlRoute := app.Group("")
	graphiqlRoute.Any("/graphiql", handleGraphiql())
}

func handleGraph(schema graphql.Schema) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "No query data"})
			return
		}

		var rBody reqBody
		err := json.NewDecoder(ctx.Request.Body).Decode(&rBody)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "Error parsing JSON request body"})
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: rBody.Query,
			Context:       context.WithValue(context.Background(), "token", ctx.Request.URL.Query().Get("token")),
		})
		if result.HasErrors() {
			errs := []string{}
			for _, err := range result.Errors {
				errs = append(errs, err.Message)
			}
			ctx.JSON(http.StatusBadRequest, result)
			return
		}
		ctx.JSON(200, result)
	}
}

func handleGraphiql() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
		if err != nil {
			fmt.Printf("failed to create new schema, error: %v", err)
		}
		graphiqlHandler.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
