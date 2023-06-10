package app

import (
	"log"
	"net/http"
	"reflect"
	"rpc/internal/database"
	"rpc/internal/services/middleware"
	"rpc/internal/transport"
	"rpc/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Server struct {
	Port       string
	Storage    storage.Storage
	Middleware middleware.Middleware
}

func NewServer(port string, storage storage.Storage, middleware middleware.Middleware) *Server {
	return &Server{
		Port:       port,
		Storage:    storage,
		Middleware: middleware,
	}
}

func (s *Server) Start() error {

	r := mux.NewRouter()

	r.HandleFunc("/login", s.Middleware.HandleLogin).Methods("POST")
	r.HandleFunc("/register", s.Middleware.HandleRegister).Methods("POST")
	r.HandleFunc("/ws", s.establishConnection)

	log.Printf("Starting server on %v", s.Port)
	return http.ListenAndServe(s.Port, r)
}

func (s *Server) establishConnection(w http.ResponseWriter, r *http.Request) {
	if !s.Middleware.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("Websocket upgrade error:", err)
		return
	}
	log.Printf("Websocket connection established from %s\n", conn.RemoteAddr().String())

	go s.websocketConnection(conn)
}

func (s *Server) websocketConnection(conn *websocket.Conn) {
	defer conn.Close()

	id := 0
	rpcHandler := transport.NewRPCHandler(s.Storage)

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

		response := result[0].Interface()
		log.Println(response)

		err := conn.WriteJSON(response)
		if err != nil {
			log.Println("error writing JSON-RPC response:", err)
			return
		}

		id++
	}
}
