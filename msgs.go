package main

import (
	"context"

	"golang.design/x/clipboard"
)

type Msg struct {
	Origin string `json:"origin"`
	Body   []byte `json:"body"`
}

func on_clip() <-chan Msg {
	c := make(chan Msg, 1)
	go func() {
		changed := clipboard.Watch(context.TODO(), clipboard.FmtText)
		for {
			message := <-changed

			c <- Msg{
				Origin: instance_id,
				Body:   message,
			}
		}
	}()

	return c
}
