package packer

import (
	"context"
	"io"
)

type PACKER interface {
	SendPackage(ctx context.Context, reader io.Writer) (*REQUEST_OBJ, RESPONSE_OBJ, error)
	ReceivePackage(ctx context.Context, writer io.Reader) (int, error)
}

type REQUEST_OBJ struct {
	Payload     []byte
	FrameType   string // e.g., "MESSAGE_FRAME", "ACK_FRAME"
	PayloadType string // e.g, JSON, Protobuf -- only protoBuf is supported for now
	PayloadSize int
}

/*
This is the response object that will be sent back to the client
It contains the connection object, stream object and the payload
*/
type RESPONSE_OBJ struct {
	Payload     []byte
	Writer      io.Writer // this is the stream writer it could be another connection or the same stream
	PayloadSize int       // this is the payload that will be sent back to the client
}
