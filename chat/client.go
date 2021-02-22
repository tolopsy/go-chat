package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	sender chan []byte
	room   *room
}

/*
reads message from the client-side and sends it to
the room's msgForwarder to be sent to all room members via their sender
channel
*/
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}

		c.room.msgForwarder <- msg
	}
}

/* receives message from client's sender channel and writes it into the socket
to be received on the client-side
*/
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.sender {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			return
		}
	}
}
