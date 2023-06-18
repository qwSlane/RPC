package app

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"rpc/internal/database"
	"rpc/internal/services/middleware"
	"rpc/internal/transport"
	"rpc/internal/transport/types"
)

type Server struct {
	Port           string
	ServiceManager *transport.ServiceManager
	Storage        storage.Storage
	Middleware     middleware.Middleware
}

func NewServer(port string, storage storage.Storage, middleware middleware.Middleware) *Server {

	serviceManager := transport.NewServiceManager()
	return &Server{
		Port:           port,
		ServiceManager: serviceManager,
		Storage:        storage,
		Middleware:     middleware,
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

	for {
		var data []byte

		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading PROTO-RPC request:", err)
			return
		}

		message := &types.Request{}

		if err := proto.Unmarshal(data, message); err != nil {
			log.Println("error unmarshaling request structure:", err)
			return

		}

		result, err := s.ServiceManager.Invoke(message.Method, message.Args)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(result)

		response, err := proto.Marshal(result)

		err = conn.WriteMessage(2, response)
		if err != nil {
			log.Println("error writing PROTO-RPC response:", err)
			return
		}

	}
}
