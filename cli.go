package main

import (
	"errors"
	"flag"
)

// FLAGS

var (
	PORT_FLAG        = flag.String("port", "", "server port")
	SERVER_HOST_FLAG = flag.String("host", "", "server host")
	PWD_FLAG         = flag.String("pass", "", "Server Password")
	CLIENT_FLAG      = flag.Bool("c", false, "run client")
	SERVER_FLAG      = flag.Bool("s", false, "run server")
	HELP_FLAG        = flag.Bool("h", false, "Show Help")
	VERBOSE_FLAG     = flag.Bool("verbose", false, "Show Verbose Logging")
)

func parse_cli_commands() error {
	flag.Parse()

	if HELP_FLAG != nil {
		if *HELP_FLAG {
			flag.PrintDefaults()
			return nil
		}
	}

	// TODO: Encrypt/Decrypt w password
	// if PWD_FLAG == nil {
	// 	return errors.New("must use password")
	// }
	//
	// if *PWD_FLAG == "" {
	// 	return errors.New("must use password")
	// }

	if CLIENT_FLAG != nil {
		if *CLIENT_FLAG {
			typ = CLIENT
		}
	}

	if SERVER_FLAG != nil {
		if *SERVER_FLAG {
			typ = SERVER
		}
	}

	if PORT_FLAG != nil {
		port = *PORT_FLAG
		if port == "" {
			return errors.New("must specify server port")
		}
	} else {
		return errors.New("must specify server port")
	}

	if typ == CLIENT {
		if SERVER_HOST_FLAG != nil {
			host = *SERVER_HOST_FLAG
			if host == "" {
				return errors.New("must specify server host")
			}
		} else {
			return errors.New("must specify server host")
		}
	}

	if VERBOSE_FLAG != nil {
		log.isVerbose = *VERBOSE_FLAG
	}

	return nil
}
