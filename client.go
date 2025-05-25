package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
)

func as_client() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s", host, port), nil)
	if err != nil {
		panic("Connection error:" + err.Error())
	}

	defer conn.Close()

	last_ping := time.Now()

	// Set up Pong Handler to respond to Pong messages
	conn.SetPongHandler(func(appData string) error {
		last_ping = time.Now()
		return nil
	})

	// Send a Ping every 5 seconds
	go func() {
		for {
			if last_ping.Before(time.Now().Add(-10 * time.Second)) {
				panic("server crashed")
			}
			err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var e Msg

			err = json.Unmarshal(msg, &e)
			if err != nil {
				continue
			}

			log.Log("client got message ", e)

			if e.Origin == instance_id {
				continue
			}

			clipboard.Write(clipboard.FmtText, e.Body)
		}
	}()

	e := on_clip()

	for {
		msg := <-e
		byts, err := json.Marshal(msg)
		if err != nil {
			continue
		}

		log.Log("client new clipboard ", string(byts))

		conn.WriteMessage(websocket.TextMessage, byts)
	}
}
