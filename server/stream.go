package server

import (
	"io"

	"github.com/admin0p/supreme-fishstick/logger"
)

func readBuffer(reader io.Reader) (*[]byte, error) {

	sizeByte := make([]byte, 1)
	_, err := reader.Read(sizeByte)
	if err != nil {
		logger.Log.Error("Failed to read size byte", "stack", err)
		return nil, err
	}

	size := int(sizeByte[0])
	payloadBuffer := make([]byte, size)
	_, err = reader.Read(payloadBuffer)
	if err != nil {
		logger.Log.Error("Failed to read payload buffer", "stack", err)
		return nil, err
	}

	return &payloadBuffer, nil
}

func sendBuffer(buffer *[]byte, writer io.Writer) {
	sizeByte := byte(len(*buffer))
	fullPackage := append([]byte{sizeByte}, *buffer...)
	_, err := writer.Write(fullPackage)
	if err != nil {
		logger.Log.Error("Failed to write buffer to writer", "stack", err)
		return
	}
	logger.Log.Info("Sent buffer to writer", "size", len(*buffer))
}
