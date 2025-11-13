package main

import (
	"io"
	"net/http"
)

type Receiver interface {
	Receive(io.Reader) (RequestMessage, error)
}

type Server struct {
	receiver Receiver
	router   *CmdRouter
}

type RequestMessage struct {
	ChatID ID
	User   User
	Text   string
}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message, err := srv.receiver.Receive(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.router.Parse(message)
}
