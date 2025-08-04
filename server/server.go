// should ideally be a quic server package
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

type StreamHandler interface {
	ProcessStream(ctx context.Context, stream *quic.Stream)
}

type QUIC_SERVER_INSTANCE struct {
	HostName string
	Port     int
	Tls      *tls.Config
	Streamer StreamHandler
	Wg       sync.WaitGroup
	Conn     *quic.Conn
}

/*
This function accepts the connections in a synchronous manner and
process the next subsequent request in non blocking fashion
*/
func (qsi *QUIC_SERVER_INSTANCE) StartServer(config *quic.Config) {

	bindAddress := qsi.generateListenAddress()
	quicListener, err := quic.ListenAddr(
		bindAddress,
		qsi.Tls,
		config,
	)

	if err != nil {
		logger.Log.Error("Failed to listen on bind address", "address", bindAddress, "stack", err)
		panic(err)
	}

	for {
		connContext := context.Background()
		conn, err := quicListener.Accept(connContext)
		if err != nil {
			logger.Log.Error("Failed to accept connection", "stack", err)
			continue
		}

		stream, err := conn.OpenStreamSync(connContext)
		if err != nil {
			logger.Log.Error("Failed to start stream", "stack", err)
			continue
		}

		//send stream connection ack
		err = sendStreamAcceptAck(stream)
		if err != nil {
			conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), "closing connection")
			continue
		}

		qsi.Streamer.ProcessStream(connContext, stream)

	}

}

func (qsi *QUIC_SERVER_INSTANCE) generateListenAddress() string {
	if qsi.HostName == "" {
		qsi.HostName = "localhost"
	}

	if qsi.Port == 0 {
		qsi.Port = 7891
	}

	return fmt.Sprintf(qsi.HostName, ":", strconv.Itoa(qsi.Port))

}

func sendStreamAcceptAck(stream *quic.Stream) error {

	ackFrame := dataframe.ACK_FRAME{
		StreamId:  int32(stream.StreamID()),
		PackId:    1,
		AckStatus: true,
	}

	ackBytes, err := proto.Marshal(&ackFrame)
	if err != nil {
		logger.Log.Error("Failed to marshal initial ack", "stack", err)
		return err
	}

	packageSize := byte(len(ackBytes))

	ackPackage := append([]byte{packageSize}, ackBytes...)

	_, err = stream.Write(ackPackage)
	if err != nil {
		logger.Log.Error("Failed to send stream accept ack", "stack", err)
		return err
	}

	return nil
}
