# Server side (Go-Cockroachdb-Graphql-Gqlgen) 

## Prerequisite

1. [Install golang](https://golang.org/dl/) [I use go 1.11 vesrion]
2. Setup GOPATH [Link1](https://golang.org/doc/code.html#GOPATH) and [Link2](https://github.com/golang/go/wiki/GOPATH)
3. [Install cockroachdb](https://www.cockroachlabs.com/docs/stable/install-cockroachdb-windows.html)
4. Check Binary file of cockroachdb : `.\cockroach.exe` run in cmd
5. Start node in cockroachdb : `cockroach start --insecure` run in cmd where cockroachdb is intalled
6. Test cluster : `cockroach sql --insecure` run in cmd where cockroachdb is intalled
7. Create database : create database dbname;(here I took chatApp)
8. Show database : show datbases;
9. Clone repo in gopath src folder (Can't run outside gopath because gqlgen not support outside gopath)
10. Install cockroachdb driver in go (go get github.com/lib/pq) for cockroachdb driver
11. Install gorilla mux (go get github.com/gorilla/mux)
12. Install gqlgen (go get -u github.com/99designs/gqlgen) [I use gqlgen 0.7.2 version]
13. If there is a new vesrion of gqlgen than do (gqlgen -v) it will generate generated.go file,resolver.go file and models_gen.go file as per new version

## Package used

1. github.com/gorilla/mux  = As router carry handler request
2. github.com/lib/pq       = For cockroachdb
3. github.com/gorilla/websocket = For handling socket 
4. github.com/99designs/gqlgen = Graphql library 
5. github.com/99designs/gqlgen/handler
6. github.com/99designs/gqlgen/graphql = For directive

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
        |--- graph
        |     |-- generated.go
        |--- model
        |     |-- models.go       
        |--- apimethods.md
        |--- Backend.md
        |--- gqlgen.yml
        |--- README.md                  
        |--- schema.graphql   
        |--- todo.txt        

## Database 

- Create tables as per configuration. 
    [You can find my database configuration file here](https://github.com/AneriShah2610/new_chat/blob/master/api/dal/db_config.json)
- Read doc to generate tables 
    [doc](https://github.com/AneriShah2610/new_chat/blob/master/dbScript/db_init.sql)
- Change configuration as per your db 

### Features

1. Create new user
2. Retrieve user
3. Live update of new user join
4. Create new chatroom
5. Retrieve all chatroom details
6. Retrieve chatroom detail by chatroomid
7. Delete chat by particular member
8. Leave chatroom only for group chat
9. Update chatroom detail
10. Add New member in chatroom
11. Delete member from chatroom
12. Live update of chatroom delete
13. Message post in chatroom
14. Live update of post message
15. Message update only by owner of message
16. Message delete only by owner of message
17. Live update of message update in group
18. Live update of message delete from group
19. Live update of new member add in chatroom
20. Retrieve member list by chatroom id
21. Update user profile
22. Remove members from chatroom only by admin

### Features 
[Api Functionality](https://github.com/AneriShah2610/new_chat/blob/master/apimethods.md)
