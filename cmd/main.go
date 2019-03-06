package main

import (
	log "log"
	http "net/http"
	os "os"
	"time"

	handler "github.com/99designs/gqlgen/handler"
	resolver "github.com/aneri/new_chat/api/handler"
	"github.com/aneri/new_chat/api/middleware"
	"github.com/aneri/new_chat/graph"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const defaultPort = "8585"

var clients = make(map[*websocket.Conn]bool)
//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(request *http.Request) bool {
//		return true
//	},
//
//}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	router := mux.NewRouter()
	//router.HandleFunc("/ws", handleConnection)
	queryHandler := corsAccess(handler.GraphQL(graph.NewExecutableSchema(
		graph.Config{
			Resolvers: &resolver.Resolver{},
		}),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(request *http.Request) bool {
				return true
			},
			HandshakeTimeout: 10 * time.Second,
		}),
	))
	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", middleware.MultipleMiddleware(queryHandler, middleware.CockroachDbMiddleware))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func corsAccess(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		response.Header().Set("Access-Control-Allow-Headers", "Accept, X-Requested-With, Content-Type, Authorization")
		next(response, request)
	})
}
//func handleConnection(w http.ResponseWriter, r *http.Request){
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil{
//		log.Println("err", err)
//	}
//	defer  ws.Close()
//	clients[ws] = true
//}