package main

import (
	"context"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

func main() {

	listener, err := quic.ListenAddr("localhost:4242", mock.GenerateDummyTLSConfig(), nil)
	if err != nil {
		logger.Log.Error("Failed to start the listener")
		return
	}
	logger.Log.Info("started listening")

	for {

		ctx := context.Background()

		conn, err := listener.Accept(ctx)
		if err != nil {
			logger.Log.Error("Failed to accept connection", "stack", err)
		}

		ConnectionHandler(conn)

		logger.Log.Info("connection accepted")

	}

}

func ConnectionHandler(conn *quic.Conn) {

	//start a stream
	stream, err := conn.OpenStreamSync(conn.Context())
	if err != nil {
		logger.Log.Error("Failed to open stream for: ", "remote_address", conn.RemoteAddr().String())
		return
	}

	// should be a protobuf
	// initial protobuf to initiate a stream
	// something like a one way handshake
	helloFrame := dataframe.MockDataFrame{
		From: "server",
		To:   conn.RemoteAddr().String(),
		Type: "Hello-init",
		Data: "Hello",
	}

	pack, err := proto.Marshal(&helloFrame)
	if err != nil {
		logger.Log.Error("Failed to serialize")
		return
	}
	packLen := byte(len(pack))
	message := append([]byte{packLen}, pack...)

	_, err = stream.Write(message)
	if err != nil {
		logger.Log.Error("Failed to initiate write to stream ")
	}

	logger.Log.Info("sent data")

}
