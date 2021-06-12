package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	sender   chan *message
	room     *room
	userData map[string]interface{}
}

/*
reads message from the client-side and sends it to
the room's msgForwarder to be sent to all room members via their sender
channel
*/
func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		//_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)

		c.room.msgForwarder <- msg
	}
}

/* receives message from client's sender channel and writes it into the socket
to be received on the client-side
*/
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.sender {
		err := c.socket.WriteJSON(msg)

		if err != nil {
			return
		}
	}
}
