package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients map[*Client]bool
	mu      sync.Mutex
}

var hub = Hub{
	clients: make(map[*Client]bool),
}

func (c *Client) readPump() {
	defer func() {
		hub.removeClient(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		log.Log("received message: ", string(message))

		var e Msg

		err = json.Unmarshal(message, &e)
		if err != nil {
			continue
		}

		if e.Origin != instance_id {
			clipboard.Write(clipboard.FmtText, e.Body)
		}

		hub.broadcast(message, c)
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func (h *Hub) broadcast(message []byte, sender *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		if client != sender {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) broadcast_all(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) addClient(c *Client) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

func (h *Hub) removeClient(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Set Ping and Pong Handlers
	conn.SetPingHandler(func(appData string) error {
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(appData string) error {
		return nil
	})

	log.Log("new connection")

	defer conn.Close()

	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.addClient(client)
	defer hub.removeClient(client)

	go client.writePump()
	client.readPump()
}

func as_server() {
	http.HandleFunc("/", wsHandler)

	go func() {
		c := on_clip()
		for {
			e := <-c
			log.Log("new clipboard item: ", e)

			byts, err := json.Marshal(e)
			if err != nil {
				continue
			}

			hub.broadcast_all(byts)
		}
	}()

	panic(http.ListenAndServe(":"+port, nil))
}
