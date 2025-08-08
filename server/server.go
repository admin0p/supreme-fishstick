// should ideally be a quic server package
package server

import (
	"context"
	"crypto/tls"
	"io"
	"strconv"
	"sync"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/admin0p/supreme-fishstick/server/packager"
	"github.com/quic-go/quic-go"
)

type ACTIVE_CLIENT_CONN map[string]*quic_conn

type SF_STREAM_HANDLER struct {
	Stream         *quic.Stream
	ActiveConn     *ACTIVE_CLIENT_CONN
	PackageHandler Packager
}
type StreamHandler interface {
	ProcessStream(ctx context.Context, streamHandler *SF_STREAM_HANDLER) error
}

type Packager interface {
	SendPackage(ctx context.Context, reader io.Writer, data any) (int, error)
	ReceivePackage(ctx context.Context, writer io.Reader, result any) (int, error)
}

type quic_conn struct {
	Conn           *quic.Conn
	UpstreamServer *QUIC_SERVER_INSTANCE
	ActiveStream   map[*quic.Stream]struct{}
}

type QUIC_SERVER_INSTANCE struct {
	HostName       string
	Port           int
	Tls            *tls.Config
	Streamer       StreamHandler
	Wg             sync.WaitGroup
	PackageEncoder Packager
	ActiveConn     ACTIVE_CLIENT_CONN
}

/*
This function accepts the connections in a synchronous manner and
process the next subsequent request in non blocking fashion
*/
func (qsi *QUIC_SERVER_INSTANCE) StartServer(config *quic.Config, packagerCode int) {

	qsi.assignServerDefaults(packagerCode)

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
	logger.Log.Info("QUIC server started", "address", bindAddress)

	for {
		connContext := context.Background()

		newConn := &quic_conn{UpstreamServer: qsi, ActiveStream: make(map[*quic.Stream]struct{})}
		newConn.Conn, err = quicListener.Accept(connContext)
		if err != nil {
			logger.Log.Error("Failed to accept connection", "stack", err)
			continue
		}

		stream, err := newConn.Conn.OpenStreamSync(connContext)
		if err != nil {
			logger.Log.Error("Failed to start stream", "stack", err)
			continue
		}

		newConn.ActiveStream[stream] = struct{}{}

		logger.Log.Info("New connection accepted", "remoteAddr", newConn.Conn.RemoteAddr().String(), "localAddr", newConn.Conn.LocalAddr().String())

		// create a new clientId ..ideally it should be done after the auth but this is a test
		clientAddr := newConn.Conn.RemoteAddr().String()
		qsi.ActiveConn[clientAddr] = newConn

		ackFrame := &dataframe.ACK_FRAME{
			StreamId:  int32(stream.StreamID()),
			PackId:    1,
			AckStatus: true,
		}
		_, err = qsi.PackageEncoder.SendPackage(connContext, stream, ackFrame)
		if err != nil {
			// this is a clean up operation for this connection
			newConn.Conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), "closing connection")
			stream.Close()
			delete(newConn.ActiveStream, stream)
			delete(qsi.ActiveConn, clientAddr)
			continue
		}

		newSfStreamHandler := SF_STREAM_HANDLER{
			Stream:         stream,
			ActiveConn:     &qsi.ActiveConn,
			PackageHandler: qsi.PackageEncoder,
		}
		//handle connection request
		for {
			// this should be a internal function that construct the request obj and passes it to the handler
			// the function should be a closure type that can take in custom handler and process request based on that
			qsi.Streamer.ProcessStream(connContext, &newSfStreamHandler)
		}

	}

}

func (qsi *QUIC_SERVER_INSTANCE) generateListenAddress() string {
	if qsi.HostName == "" {
		qsi.HostName = "localhost"
	}

	if qsi.Port == 0 {
		qsi.Port = 7891
	}

	return qsi.HostName + ":" + strconv.Itoa(qsi.Port)

}

func (qsi *QUIC_SERVER_INSTANCE) assignServerDefaults(packagerCode int) {
	if qsi.Tls == nil {
		qsi.Tls = mock.GenerateDummyTLSConfig()
	}

	if qsi.ActiveConn == nil {
		qsi.ActiveConn = make(ACTIVE_CLIENT_CONN)
	}

	if qsi.Streamer == nil {
		qsi.Streamer = &ProtoHandler{}
	}

	if packagerCode == 1 {
		qsi.PackageEncoder = &packager.PROTO_ENCODE{}
	}

}
