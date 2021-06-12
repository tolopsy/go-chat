package main

import (
	"go-projects/go-chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	addMember    chan *client
	removeMember chan *client
	roomMembers  map[*client]bool
	msgForwarder chan *message
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
				r.tracer.Trace("Message sent: ", newMessage.Message)
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
		msgForwarder: make(chan *message),
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

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie: ", err)
		return
	}

	newClient := &client{
		socket:   socket,
		sender:   make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.addMember <- newClient

	defer func() { r.removeMember <- newClient }()
	go newClient.write()
	newClient.read()
}
