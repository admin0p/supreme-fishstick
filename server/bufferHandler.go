package server

import (
	"io"

	"github.com/admin0p/supreme-fishstick/logger"
)

const (
	MAX_BUFFER_SIZE = 10 << 20
)

func readBuffer(reader io.Reader) (*[]byte, error) {
	sizeByte := make([]byte, 1)
	_, err := reader.Read(sizeByte)
	if err != nil {
		return nil, NewErr(err, BUFF_READ_FAILED, "Failed to read size buffer")
	}

	size := int(sizeByte[0])
	if size > MAX_BUFFER_SIZE {
		return nil, NewErr(nil, BUFF_SIZE_EXCEEDED, "exceeded read buffer size")
	}

	payloadBuffer := make([]byte, size)
	_, err = reader.Read(payloadBuffer)
	if err != nil {
		return nil, NewErr(err, BUFF_READ_FAILED, "Failed to read payload buffer")
	}
	return &payloadBuffer, nil
}

func sendBuffer(buffer *[]byte, writer io.Writer) error {
	size := len(*buffer)
	if size > MAX_BUFFER_SIZE {
		return NewErr(nil, BUFF_SIZE_EXCEEDED, "exceeded read buffer size")
	}

	sizeByte := byte(size)
	fullPackage := append([]byte{sizeByte}, *buffer...)
	_, err := writer.Write(fullPackage)
	if err != nil {
		return NewErr(err, BUFF_WRITE_FAILED, "Failed to write buffer to writer")
	}

	logger.Log.Info("Sent buffer to writer", "size", size)
	return nil
}
