package main

import "share/server"

func main() {
	server, err := server.NewServer()
	if err != nil {
		return
	}

	server.Start()
}
