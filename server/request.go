package server

import (
	"context"

	"github.com/admin0p/supreme-fishstick/logger"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

// look for better name
type SF_REQ_CONTEXT struct {
	Ctx            context.Context
	ProtoPkg       proto.Message
	Payload        *[]byte
	upstreamServer *QUIC_SERVER_INSTANCE
}

func (sr *SF_REQ_CONTEXT) GetConnectionStream(connId string, streamName string) *quic.Stream {
	conn, connExist := sr.upstreamServer.ActiveConn[connId]
	if !connExist {
		logger.Log.Error("Conn does not exists")
		return nil
	}

	stream, streamExist := conn.ActiveStream[streamName]
	if !streamExist {
		logger.Log.Error("Stream does not exists")
		return nil
	}

	return stream
}
