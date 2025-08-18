package server

import (
	"context"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/server/proto"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

type QuicConn struct {
	Conn           *quic.Conn
	UpstreamServer *QUIC_SERVER_INSTANCE
	ActiveStream   ACTIVE_STREAM
}

// this function should be a closure type that can take in custom handler and process request based on that
// this is for a single connection only
// should be run in a goroutine
func (qc *QuicConn) serve(ctx context.Context) {

	defaultStream := qc.ActiveStream["default"]
	if defaultStream == nil {
		logger.Log.Error("No default stream found for the connection")
		return
	}

	// TODO: check for quit signal
	for {

		rawPayload, err := readBuffer(defaultStream)
		if err != nil {
			logger.Log.Error("Failed to read buffer from stream", "stack", err)
			break
		}

		messageFrame := &dataframe.MESSAGE_FRAME{}
		err = proto.Unmarshal(*rawPayload, messageFrame)
		if err != nil {
			logger.Log.Error("Failed to unmarshal protobuf message", "stack", err)
			break
		}

		logger.Log.Info("Received message frame", "message", messageFrame.Payload)
		// process using handler

		// this will ideally be a ref from qsi or upstream server
		err = RequestHandler(messageFrame, defaultStream)
		if err != nil {
			logger.Log.Error("Failed to handle request", "stack", err)
		}

	}

	return
}

func (qc *QuicConn) GetStreamFromClientId(clientId string) (*quic.Stream, error) {

	clientConn, ok := qc.UpstreamServer.ActiveConn[clientId]
	if !ok {
		logger.Log.Error("Client connection not found", "clientId", clientId)
		return nil, nil
	}
	stream, ok := clientConn.ActiveStream["default"]
	if !ok {
		logger.Log.Error("Stream not found for client", "clientId", clientId)
		return nil, nil
	}

	return stream, nil
}

func RequestHandler(request *dataframe.MESSAGE_FRAME, stream *quic.Stream) error {
	logger.Log.Info("Handling request", "message", request.GetPayload())

	// Here you can implement your logic to handle the request
	// For example, you might want to send an ACK response back

	ackFrame := &dataframe.ACK_FRAME{
		PackId:    request.GetMessageId() + 1,
		StreamId:  int32(request.GetStreamId()),
		AckStatus: true,
	}
	rawAckBuffer, err := proto.Marshal(ackFrame)
	if err != nil {
		logger.Log.Error("Failed to marshal ACK frame", "stack", err)
		return err
	}

	sendBuffer(&rawAckBuffer, stream)
	return nil

}
