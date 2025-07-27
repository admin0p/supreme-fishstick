package main

import (
	"context"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"github.com/admin0p/supreme-fishstick/serializer"
	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/quic-go/quic-go"
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

	ctx := context.Background()
	//start a stream
	stream, err := StartStream(ctx, conn)
	if err != nil {
		logger.Log.Error("Failed to start the stream", "stack", err)
		return
	}
	defer stream.Close()

	for {

		messageFrame := dataframe.MESSAGE_FRAME{}
		err = serializer.DeserializePackage(ctx, stream, &messageFrame)
		if err != nil {
			logger.Log.Error("Failed to deserialize", "stack", err)
			return
		}
		logger.Log.Info(
			"received ===> ",
			"from", messageFrame.GetFrom(),
			"to", messageFrame.GetTo(),
			"type", messageFrame.GetType(),
			"message_format", messageFrame.GetMessageFormat(),
			"msg", messageFrame.GetMessage(),
		)

		if messageFrame.GetMessage() == "quit\n" {
			logger.Log.Info("received quit signal")
			return
		}

		ackFrame := dataframe.STREAM_HELLO{Message: "ack"}
		err = serializer.SerializePayloadAndSend(ctx, stream, &ackFrame)
		if err != nil {
			logger.Log.Error("Failed to send ack", "stack", err)
			return
		}
	}

}

func StartStream(ctx context.Context, conn *quic.Conn) (*quic.Stream, error) {

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		logger.Log.Error("Failed to start a stream: ")
		return nil, err
	}

	baseHelloFrame := dataframe.STREAM_HELLO{
		StreamId: int64(stream.StreamID()),
		Message:  "Stream hello",
	}

	err = serializer.SerializePayloadAndSend(ctx, stream, &baseHelloFrame)
	if err != nil {
		logger.Log.Error("Failed to serialize and send ", "stack", err)
		return nil, err
	}

	return stream, nil
}
