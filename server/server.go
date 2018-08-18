package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/emit-sh/emit/server/storage"
	"github.com/gorilla/mux"
	"github.com/kurin/blazer/b2"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Port int

	tlsConfig *tls.Config
	storage   storage.Storage
	email     EmailSender
}

type OptionFn func(*Server)

func NewServer() (server *Server, err error) {
	//TODO: setup server, maybe just from env variables? cmd line?
	server = new(Server)
	server.storage, err = storage.NewDigitalOceanStorage()
	server.email = createMailgunEmailClient()

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
	/*
	mgKey := os.Getenv("B2_ACCOUNT_KEY")
	mgPubKey := os.Getenv("B2_ACCOUNT_PUB_KEY")
	mgDomain := os.Getenv("MG_DOMAIN")
	*/
	mail = newMailGunSender("mg.emit.sh")
	return
}

func (server Server) SetPort(port int) {
	server.Port = port
}

func (server Server) Start() {

	r := mux.NewRouter()
	r.HandleFunc("/", server.HomePage).Methods("GET")
	r.HandleFunc("/", server.FileHandler).Methods("POST")
	r.HandleFunc("/{token}/{filename}", server.Download).
		Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
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
