package serializer

import (
	"context"
	"io"

	"github.com/admin0p/supreme-fishstick/logger"
	"google.golang.org/protobuf/proto"
)

func Receive(ctx context.Context, reader io.Reader, result proto.Message) error {

	sizeByte := make([]byte, 1)

	sizeRead, err := reader.Read(sizeByte)
	if err != nil {
		logger.Log.Error("Failed to read the size byte")
		return err
	}

	logger.Log.Info("Read the size byte: ", "byteSize", sizeRead, "size", sizeByte)

	payloadBuffer := make([]byte, sizeByte[0])

	_, err = reader.Read(payloadBuffer)
	if err != nil {
		logger.Log.Error("Failed to REad payload")
		return err
	}

	proto.Unmarshal(payloadBuffer, result)
	return nil
}

func Send(ctx context.Context, writer io.Writer, payload proto.Message) error {

	payloadByte, err := proto.Marshal(payload)
	if err != nil {
		logger.Log.Error("Failed to marshall to proto")
		return err
	}

	payloadLen := byte(len(payloadByte))

	packageFrame := append([]byte{payloadLen}, payloadByte...)

	_, err = writer.Write(packageFrame)
	if err != nil {
		logger.Log.Error("Failed to send the package")
		return err
	}

	return nil
}
