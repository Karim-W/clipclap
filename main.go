package main

import (
	"flag"
	"fmt"
	"os"

	"golang.design/x/clipboard"
)

type (
	MODE int64
)

var (
	port        string
	host        string
	password    string
	typ         MODE
	instance_id string
)

const (
	HELP MODE = iota
	SERVER
	CLIENT
)

func recoverPanic() {
	// This function will recover from a panic
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}

func main() {
	log = Logger{}

	hostname, err := os.Hostname()
	assert(err)

	instance_id = hostname

	err = parse_cli_commands()
	assert(err)
	err = clipboard.Init()
	if err != nil {
		panic(err)
	}

	switch typ {
	case SERVER:
		as_server()
	case CLIENT:
		defer recoverPanic()
		as_client()
	default:
		flag.PrintDefaults()
		return
	}
}

func assert(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
