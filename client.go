package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
)

func as_client() {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s", host, port), nil)
	if err != nil {
		panic("Connection error:" + err.Error())
	}
	defer conn.Close()

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

		conn.WriteMessage(websocket.TextMessage, byts)
	}
}
