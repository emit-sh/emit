package server

import (
	"crypto/tls"
	"github.com/kurin/blazer/b2"
	"os"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"emit/server/storage"
	"fmt"
)

type Server struct {
	Port int

	tlsConfig *tls.Config
	storage   storage.Storage

}

type OptionFn func(*Server)

func NewServer() (server *Server, err error) {
	//TODO: setup server, maybe just from env variables? cmd line?
	server = new(Server)
	server.storage, err = createBackBlazeStorage()

	if err != nil {
		fmt.Print(err)
	}
	return
}

func createBackBlazeStorage() (s storage.BackBlazeStorage, err error) {
	b2id := os.Getenv("B2_ACCOUNT_ID")
	b2key := os.Getenv("B2_ACCOUNT_KEY")
	s = storage.BackBlazeStorage{}
	ctx := context.Background()
	s.Client, err = b2.NewClient(ctx, b2id, b2key)
	return
}

func createMailgunEmailClient() (mail MailgunEmailSender) {
	mgKey := os.Getenv("B2_ACCOUNT_KEY")
	mgPubKey := os.Getenv("B2_ACCOUNT_PUB_KEY")
	mgDomain := os.Getenv("MG_DOMAIN")
	mail = newMailGunSender(mgDomain,mgKey,mgPubKey)
	return
}

func (server Server) SetPort(port int) {
	server.Port = port
}

func (server Server) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/", server.FileHandler).Methods("POST")
	r.HandleFunc("/{token}/{filename}", server.Download).
		Methods("GET")
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