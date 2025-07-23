package main

import (
	"context"
	"crypto/tls"
	"fmt"

	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
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

	sizeByte := make([]byte, 1)

	_, err = newStream.Read(sizeByte)
	if err != nil {
		fmt.Println("Failed to read message")
		return
	}

	readBuffer := make([]byte, int(sizeByte[0]))
	_, err = newStream.Read(readBuffer)
	if err != nil {
		fmt.Println("Failed to read message")
		return
	}

	message := dataframe.MockDataFrame{}

	err = proto.Unmarshal(readBuffer, &message)
	if err != nil {
		fmt.Println("Failed to read the proto buffer")
		return
	}

	// fmt.Println("readBytes ==> ", readBytes)
	fmt.Println("data ==?> ", message.GetData())
}
