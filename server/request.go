package server

import (
	"context"

	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

// look for better name
type SF_REQ_CONTEXT struct {
	Ctx            context.Context
	ProtoPkg       proto.Message
	CurrentStream  *quic.Stream
	Payload        *[]byte
	upstreamServer *QUIC_SERVER_INSTANCE
}

func (sr *SF_REQ_CONTEXT) getConnectionStream(connId string, streamName string) (*quic.Stream, error) {
	conn, connExist := sr.upstreamServer.ActiveConn[connId]
	if !connExist {
		return nil, NewErr(nil, INVALID_CONN, "Connection does not exists in this server")
	}

	stream, streamExist := conn.ActiveStream[streamName]
	if !streamExist {
		return nil, NewErr(nil, INVALID_STREAM_REF, "Failed to fetch stream "+streamName)
	}

	return stream, nil
}

// send on current stream
func (sr *SF_REQ_CONTEXT) SendOnSameStream(rawMsg *[]byte, streamName string) error {

	if sr.CurrentStream == nil {
		// return new error
		return NewErr(nil, INVALID_STREAM_REF, "Failed to fetch stream "+streamName)

	}

	// check for error this function should return an error
	err := sendBuffer(rawMsg, sr.CurrentStream)
	if err != nil {
		return err
	}

	return nil
}

// send to other connection stream connected on the same server
func (sr *SF_REQ_CONTEXT) Send(rawMsg *[]byte, toConn string, streamType string) {
	// validate the to string - xyz@domain.ext
	// check if its the same domain
	// if not same domain send to the registered federated server
}
