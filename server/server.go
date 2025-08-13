// should ideally be a quic server package
package server

import (
	"context"
	"crypto/tls"
	"io"
	"strconv"
	"sync"

	"github.com/admin0p/supreme-fishstick/logger"
	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/admin0p/supreme-fishstick/server/packer"
	"github.com/quic-go/quic-go"
)

type ACTIVE_CLIENT_CONN map[string]*QuicConn
type ACTIVE_STREAM map[string]*quic.Stream

// This is not required as it is only used in processing a stream which is a method of QUIC_SERVER_INSTANCE so we can use that instead

type QuicConn struct {
	Conn           *quic.Conn
	UpstreamServer *QUIC_SERVER_INSTANCE
	ActiveStream   ACTIVE_STREAM
}

type QUIC_SERVER_INSTANCE struct {
	HostName       string
	Port           int
	Tls            *tls.Config
	Handler        HANDLER
	Wg             sync.WaitGroup
	PackageEncoder packer.PACKER
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

		newConn := &QuicConn{UpstreamServer: qsi, ActiveStream: make(ACTIVE_STREAM)}
		newConn.Conn, err = quicListener.Accept(connContext)
		if err != nil {
			logger.Log.Error("Failed to accept connection", "stack", err)
			continue
		}

		stream, err := newConn.Conn.AcceptStream(connContext)
		if err != nil {
			logger.Log.Error("Failed to start stream", "stack", err)
			continue
		}
		newConn.ActiveStream["default"] = stream

		// only when a stream has started we can consider the connection as active
		clientAddr := newConn.Conn.RemoteAddr().String()
		qsi.ActiveConn[clientAddr] = newConn
		logger.Log.Info("New connection accepted", "remoteAddr", newConn.Conn.RemoteAddr().String(), "localAddr", newConn.Conn.LocalAddr().String())

		//handle connection request
		for {
			// this should be a internal function that construct the request obj and passes it to the handler
			// the function should be a closure type that can take in custom handler and process request based on that
			//qsi.Streamer.ProcessStream(connContext, &newSfStreamHandler)
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

	// if qsi.Streamer == nil {
	// 	qsi.Streamer = &ProtoHandler{}
	// }

	// if packagerCode == 1 {
	// 	qsi.PackageEncoder = &packager.PROTO_ENCODE{}
	// }

}

type HANDLER func(ctx context.Context, request *packer.REQUEST_OBJ, resp packer.RESPONSE_OBJ) error

// this function should be a closure type that can take in custom handler and process request based on that
func (qc *QuicConn) serve(ctx context.Context) error {
	//TODO:  handle the nil case and assign default handler
	// CREATE a REQUEST OBJECT and a RESPONSE object

	//1 . get default stream from the connection
	defaultStream := qc.ActiveStream["default"]
	if defaultStream == nil {
		logger.Log.Error("No default stream found for the connection")
		return nil
	}

	go func() {
		for {

			// read
			sizeBuffer := make([]byte, 1) // Adjust size as needed
			_, err := defaultStream.Read(sizeBuffer)
			if err != nil && err != io.EOF {
				logger.Log.Error("Failed to read from stream", "stack", err)
			}

		}
	}()

	return nil
}
