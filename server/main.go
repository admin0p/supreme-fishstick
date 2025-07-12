package main

import (
	"context"
	"fmt"

	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/quic-go/quic-go"
)

func main() {

	l, err := quic.ListenAddr("localhost:4242", mock.GenerateDummyTLSConfig(), nil)
	if err != nil {
		fmt.Println("Error starting QUIC listener : ", err)
		return
	}

	for {
		ctx := context.Background()
		ses, err := l.Accept(ctx)
		if err != nil {
			fmt.Println("Failed to create session : ", err)
		}
		fmt.Print("accepted connection \n")

		ReqHandler(ses)

	}

}

func ReqHandler(ses *quic.Conn) {
	requestContext := context.Background()
	stream, err := ses.AcceptStream(requestContext)
	if err != nil {
		fmt.Println("failed to accept stream == ", err)
		return
	}
	defer stream.Close()
	fmt.Println("started a new incoming stream")

}
