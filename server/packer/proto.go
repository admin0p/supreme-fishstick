package packer

import (
	"context"
	"io"
)

type PROTO_PACKER struct {
}

func (PROTO_PACKER) ReceivePackage(ctx context.Context, reader io.Reader) (int, error) {
	// Implement the logic to receive a package using Protobuf
	// This is a placeholder implementation
	// y := PROTO_PACKER{}
	// z := reflect.New(reflect.TypeOf(y).Elem())

	return 0, nil
}
