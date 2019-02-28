# Chat Prototype
> Chat prototype using Golang

## Prerequisite

1. [Install golang](https://golang.org/dl/)[Here, code base on go 1.11 version]
2. Setup GOPATH [Link1](https://golang.org/doc/code.html#GOPATH) and [Link2](https://github.com/golang/go/wiki/GOPATH)
3. [Install cockroachdb](https://www.cockroachlabs.com/docs/stable/install-cockroachdb-windows.html)
4. Check Binary file of cockroachdb : `.\cockroach.exe` run in cmd
5. Start node in cockroachdb : `cockroach start --insecure` run in cmd where cockroachdb is intalled
6. [Install Glide](https://github.com/Masterminds/glide)

## Package used

1. github.com/gorilla/mux  = As router carry handler request
2. github.com/lib/pq       = For cockroachdb
3. github.com/gorilla/websocket = For handling socket 
4. github.com/99designs/gqlgen = Graphql library 
5. github.com/99designs/gqlgen/handler
6. github.com/99designs/gqlgen/graphql = For directive

## Getting Stared

1. Clone the repo under `$GOPATH/src/zuru.tech`.If does nt exist than first create it.Then run `git clone https://github.com/AneriShah2610/new_chat.git`
2. Run `glide install`(gqlgen version 0.7.2)
3. Test cockroachdb cluster: `cockroach sql --insecure` run in cmd where cockroachdb is insalled
4. Create Database: `create database chat_tets if not exist`
5. Add dbScript in database.Find [docs](https://github.com/AneriShah2610/new_chat/tree/master/dbScript) for dbScript
6. Run `go run cmd/mian.go`
7. Open `https://localhost:8585/` for GraphQL Playground


## Folder structure of server side

        |--- api
        |     |--- dal
        |     |     |-- db_config.json        
        |     |     |-- connection.go       
        |     |     |-- db_config.go         
        |     |--- directive
        |     |     |-- directive.go        
        |     |--- handler
        |     |     |-- resolver.go        
        |     |     |-- user_handler.go
        |     |     |-- chatroom_handler.go
        |     |     |-- members_handler.go
        |     |     |-- chatconversation_handler.go
        |     |--- helper
        |     |     |-- helpers.go
        |     |--- middleware
        |     |     |-- middlewares.go
        |--- cmd
        |     |-- main.go
        |--- dbScript
        |     |-- db_init.sql
        |--- error
        |     |- error.go
        |--- graph
        |     |-- generated.go
        |--- model
        |     |-- models.go       
        |--- apimethods.md
        |--- glide.lock
        |--- glide.yml
        |--- gqlgen.yml
        |--- README.md                  
        |--- schema.graphql 
        |--- schema.md  
        |--- todo.txt      
        |--- WholechatReadme.md  

## Database 

- Create tables as per configuration. 
    [You can find my database configuration file here](https://github.com/AneriShah2610/new_chat/blob/master/api/dal/db_config.json)
- Read doc to generate tables 
    [doc](https://github.com/AneriShah2610/new_chat/blob/master/dbScript/db_init.sql)
- Change configuration as per your db 

### Chat Prototype API features 
- [Api Functionality](https://github.com/AneriShah2610/new_chat/blob/master/apimethods.md)

#### Add API features

1. Alter desire changes in `schema.graphql` file and create appropriate models under `api/model` package
2. Also update `gqlgen.yml` file to add new model mapping
3. Run `gqlgen`
4. Implement new resolvers in `api/handler` package 

#### Schema
- [Schema Docs](https://github.com/AneriShah2610/new_chat/blob/master/schema.md) - generated using [graphql-markdown](https://www.npmjs.com/package/graphql-markdown)
