package packager

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

var MAX_MESSAGE_SIZE = 1 << 20 // 1 MB

type PROTO_ENCODE struct {
	Message *proto.Message
}

func (p *PROTO_ENCODE) ReceivePackage(ctx context.Context, reader io.Reader, result any) (int, error) {
	sizeByte := make([]byte, 1)
	_, err := reader.Read(sizeByte)
	if err != nil {
		return 0, err
	}

	// if int(sizeByte[0]) > MAX_MESSAGE_SIZE {
	// 	return 0, io.ErrShortBuffer
	// }

	payloadBuffer := make([]byte, sizeByte[0])
	nBytesRead, err := reader.Read(payloadBuffer)
	if err != nil {
		return nBytesRead, err
	}

	err = proto.Unmarshal(payloadBuffer, result.(proto.Message))
	if err != nil {
		return nBytesRead, err
	}

	result = result.(proto.Message)

	return nBytesRead, nil
}

func (p *PROTO_ENCODE) SendPackage(ctx context.Context, writer io.Writer, data any) (int, error) {

	protoData, ok := data.(proto.Message)
	if !ok {
		// change the error
		fmt.Println("Data is not a proto message")
		return 0, io.ErrUnexpectedEOF
	}

	payloadByte, err := proto.Marshal(protoData)
	if err != nil {
		fmt.Println("Failed to marshal proto message:", err)
		return 0, err
	}

	sizeByte := byte(len(payloadByte))
	fmt.Println("Payload size:", sizeByte)
	protoPackage := append([]byte{sizeByte}, payloadByte...)
	fmt.Println("Proto package size:", protoPackage)
	nBytesWritten, err := writer.Write(protoPackage)
	if err != nil {
		fmt.Println("Failed to write proto package:", err)
		return nBytesWritten, err
	}
	fmt.Println("Wrote proto package of size:", nBytesWritten)

	return nBytesWritten, nil
}
