package server

import (
	"crypto/tls"
	"github.com/kurin/blazer/b2"
	"os"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Server struct {
	Port int

	tlsConfig *tls.Config
	storage   Storage
}

type OptionFn func(*Server)

func NewServer() (server *Server, err error) {
	//TODO: setup server, maybe just from env variables? cmd line?
	server = new(Server)
	server.storage, err = createBackBlazeStorage()
	return
}

func createBackBlazeStorage() (storage BackBlazeStorage, err error) {
	b2id := os.Getenv("B2_ACCOUNT_ID")
	b2key := os.Getenv("B2_ACCOUNT_KEY")
	storage = BackBlazeStorage{}
	ctx := context.Background()
	storage.Client, err = b2.NewClient(ctx, b2id, b2key)
	return
}

func (server Server) SetPort(port int) {
	server.Port = port
}

func (server Server) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/", server.FileHandler).Methods("POST")
	r.HandleFunc("/{token}/{filename}", server.Download).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func Port(i int) OptionFn {
	return func(server *Server) {
		server.Port = i
	}
}