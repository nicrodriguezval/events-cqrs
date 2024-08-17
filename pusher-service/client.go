package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	socket   *websocket.Conn
	outbound chan []byte
	id       string
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (c *Client) Write() error {
	for message := range c.outbound {
		err := c.socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
	}

	// Handle closed channel gracefully
	err := c.socket.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		return fmt.Errorf("failed to send close message: %w", err)
	}

	return nil
}

func (c Client) Close() {
	c.socket.Close()
	close(c.outbound)
}
