package main

import (
	"go-chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	addMember    chan *client
	removeMember chan *client
	roomMembers  map[*client]bool
	msgForwarder chan []byte
	tracer       trace.Tracer
}

// handles the main activities in the room (join, leave and forward messages)
func (r *room) run() {
	for {
		select {
		case newMember := <-r.addMember:
			r.roomMembers[newMember] = true
			r.tracer.Trace("New Client Joined")

		case exMember := <-r.removeMember:
			delete(r.roomMembers, exMember)
			close(exMember.sender)
			r.tracer.Trace("Client left")

		case newMessage := <-r.msgForwarder:
			for mem := range r.roomMembers {
				mem.sender <- newMessage
				r.tracer.Trace("-- sent to client")
			}
		}

	}
}

// createNewRoom returns a new room object
func createNewRoom() *room {
	return &room{
		addMember:    make(chan *client),
		removeMember: make(chan *client),
		roomMembers:  make(map[*client]bool),
		msgForwarder: make(chan []byte),
		tracer:       trace.Off(),
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 512
)

// upgrader will upgrade http connections to websocket connection
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil) // upgrades http connection to websocket connection
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	newClient := &client{
		socket: socket,
		sender: make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.addMember <- newClient

	defer func() { r.removeMember <- newClient }()
	go newClient.write()
	newClient.read()
}
