package main

import (
	"fmt"
	"io"
	"net/http"
)

type Receiver interface {
	Receive(io.Reader) (RequestMessage, error)
}

type Server struct {
	receiver Receiver
}

type RequestMessage struct {
	ChatID int
	Text   string
}

func (srv Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	message, err := srv.receiver.Receive(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	}

	fmt.Fprintf(w, "message: %+v", message)
}
