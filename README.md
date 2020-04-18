# Restful GIN
Rest API with Golang, MongoDB

# Feature
* CRUD API
* Authentication
* Authorization
* CORS
* Auto generate Swagger Docs with annotations

# Technologies
* [Gin](https://github.com/gin-gonic/gin)
* [MongoDB](https://www.mongodb.com)

# Set up
* Create file .env
* Set MongoDB URI and DB
  - PORT = "8585" or your port
  - MONGO_HOST = "your host/ localhost:27017"
  - MONGO_DB_NAME = "your db name"

# Run
* `go mod download` for download dependencies
* `go run main.go`

# Swagger
* `localhost:8585/swagger/index.html`

# Graphql
* `localhost:8585/graphql`
* `localhost:8585/graphiql` (UI)

