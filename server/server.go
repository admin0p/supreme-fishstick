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
	"github.com/admin0p/supreme-fishstick/server/packer"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
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
		go func() {
			for {
				newConn.serve(connContext)
			}
		}()

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
// this is for a single connection only
// should be run in a goroutine
func (qc *QuicConn) serve(ctx context.Context) {

	defaultStream := qc.ActiveStream["default"]
	if defaultStream == nil {
		logger.Log.Error("No default stream found for the connection")
		return
	}

	// TODO: check for quit signal
	for {

		rawPayload, err := readBuffer(defaultStream)
		if err != nil {
			logger.Log.Error("Failed to read buffer from stream", "stack", err)
			break
		}

		messageFrame := &dataframe.MESSAGE_FRAME{}
		err = proto.Unmarshal(*rawPayload, messageFrame)
		if err != nil {
			logger.Log.Error("Failed to unmarshal protobuf message", "stack", err)
			break
		}

		logger.Log.Info("Received message frame", "message", messageFrame.GetMessage())
		// process using handler
		err = RequestHandler(messageFrame, defaultStream)
		if err != nil {
			logger.Log.Error("Failed to handle request", "stack", err)
		}

	}

	return
}

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

func RequestHandler(request *dataframe.MESSAGE_FRAME, stream *quic.Stream) error {
	logger.Log.Info("Handling request", "message", request.GetMessage())

	// Here you can implement your logic to handle the request
	// For example, you might want to send an ACK response back

	ackFrame := &dataframe.ACK_FRAME{
		PackId:    request.GetPackId() + 1,
		StreamId:  int32(request.GetStreamId()),
		AckStatus: true,
	}
	rawAckBuffer, err := proto.Marshal(ackFrame)
	if err != nil {
		logger.Log.Error("Failed to marshal ACK frame", "stack", err)
		return err
	}

	sendBuffer(&rawAckBuffer, stream)
	return nil

}
