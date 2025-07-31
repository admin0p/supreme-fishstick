package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"github.com/admin0p/supreme-fishstick/serializer"
	"github.com/quic-go/quic-go"
)

/*
 ** TODOS:
 1. Create a logger package
 2. Use logger instead of fmt.Println
 3. Make sure all TLS related configs come from a config file
 4. Server opens the stream after connection
*/

func main() {

	ctx := context.Background()
	// should be handled via a config
	c, err := quic.DialAddr(ctx, "localhost:4242", &tls.Config{InsecureSkipVerify: true, NextProtos: []string{"quick"}}, nil)
	if err != nil {
		fmt.Println("Error connecting to server : ", err)
		return
	}
	fmt.Println("connection accepted")

	newStream, err := c.AcceptStream(ctx)
	if err != nil {
		fmt.Println("Failed to Accept stream")
		return
	}
	defer newStream.Close()

	message := dataframe.ACK_FRAME{}

	err = serializer.Receive(ctx, newStream, &message)
	if err != nil {
		logger.Log.Error("Failed to deserialize package", "stack", err)
		return
	}
	fmt.Println("received Message ", message.GetStreamId(), message.GetAckStatus(), message.GetPackId())

	for {
		input := ReadInput()
		if input == "quit" {
			break
		}

		payload := dataframe.MESSAGE_FRAME{
			Message:       input,
			From:          c.LocalAddr().String(),
			To:            c.RemoteAddr().String(),
			MessageFormat: "string",
			Type:          "message",
			PackId:        message.GetPackId() + 1,
		}

		err = serializer.Send(ctx, newStream, &payload)
		if err != nil {
			fmt.Println(err)
			break
		}

		ackMessage := dataframe.ACK_FRAME{}

		err = serializer.Receive(ctx, newStream, &ackMessage)

		if err != nil {
			if err == io.EOF {
				logger.Log.Info("Stream closed")
				return
			}
			logger.Log.Error("Failed to deserialize package", "stack", err)
			return
		}

		if !ackMessage.GetAckStatus() {
			fmt.Println("message delivery failed not acked")
			return
		}

	}
}

func ReadInput() string {

	fmt.Println("Enter you message")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return input
}
