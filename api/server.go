package api

import (
	"log"
	"main/storage"
	"main/types"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	Port       string
	Storage    storage.Storage
	Middleware Middleware
}

func NewServer(port string, storage storage.Storage, middleware Middleware) *Server {
	return &Server{
		Port:       port,
		Storage:    storage,
		Middleware: middleware,
	}
}

func (s *Server) Start() error {

	r := mux.NewRouter()

	r.HandleFunc("/login", s.Middleware.handleLogin).Methods("POST")
	r.HandleFunc("/register", s.Middleware.handleRegister).Methods("POST")
	r.HandleFunc("/ws", s.establishConnection)

	log.Printf("Starting server on %v", s.Port)
	return http.ListenAndServe(s.Port, r)
}

func (s *Server) establishConnection(w http.ResponseWriter, r *http.Request) {
	if !s.Middleware.isAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}
	log.Printf("Websocket connection established from %s\n", conn.RemoteAddr().String())

	go websocketConnection(conn)
}

func websocketConnection(conn *websocket.Conn) {
	defer conn.Close()

	id := 0
	rpcHandler := &Handler{}

	for {
		var request types.RpcRequest
		if err := conn.ReadJSON(&request); err != nil {
			log.Println("error reading JSON-RPC request:", err)
			return
		}

		argsValue := make([]reflect.Value, len(request.Args))
		for i, arg := range request.Args {
			argsValue[i] = reflect.ValueOf(arg)
		}
		result := reflect.ValueOf(rpcHandler).MethodByName(request.Method).Call(argsValue)
		if result != nil {
			response := types.RpcResponse{Result: result[0].Interface(), Error: "", Id: request.Id}

			if err := conn.WriteJSON(response); err != nil {
				log.Println("error writing JSON-RPC response:", err)
				return
			}
		}

		id++
	}
}
