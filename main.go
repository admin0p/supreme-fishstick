package main

import "github.com/admin0p/supreme-fishstick/server"

func main() {

	quicServer := &server.QUIC_SERVER_INSTANCE{
		HostName: "localhost",
		Port:     4242,
	}

	quicServer.StartServer(nil, 1)
}
