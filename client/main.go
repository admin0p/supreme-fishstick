package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

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

	newStream, err := c.OpenStream()
	if err != nil {
		fmt.Println("Failed to open stream")
		return
	}
	newStream.Write([]byte("hello"))
	// newStream.CancelWrite(quic.StreamErrorCode(quic.NoError))
	time.Sleep(1 * time.Second)
	fmt.Println(" \n opened new stream")
	newStream.Close()
	// receive message from server
	// buff := make([]byte, 1024)
	// readIndex, err := newStream.Read(buff)
	// if err != nil {
	// 	fmt.Println("read failed ==<> ", err)
	// 	return
	// }
	// fmt.Println("ack message received ==> ", readIndex, " message <> ", string(buff))
}
