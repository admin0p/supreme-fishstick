// should ideally be a quic server package
package server

import (
	"context"
	"crypto/tls"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/admin0p/supreme-fishstick/logger"
	mock "github.com/admin0p/supreme-fishstick/server/mocks"
	"github.com/quic-go/quic-go"
)

type ACTIVE_CLIENT_CONN map[string]*QuicConn
type ACTIVE_STREAM map[string]*quic.Stream
type HANDLER interface {
	AuthHandler(SR *SF_REQ_CONTEXT) (bool, error)
	messageHandler(SR *SF_REQ_CONTEXT) error
}

type QUIC_SERVER_INSTANCE struct {
	HostName   string
	Port       int
	Tls        *tls.Config
	Handler    HANDLER
	Wg         sync.WaitGroup
	ActiveConn ACTIVE_CLIENT_CONN
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

	// check for termination
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)
	quitSig := make(chan bool, 1)
	go func() {
		sig := <-osSignal
		logger.Log.Info("Received Quit signal, quitting gracefully ...", "signal", sig)
		quitSig <- true
	}()

	for !<-quitSig {
		baseCtx := context.Background()
		newConn, err := quicListener.Accept(baseCtx)
		if err != nil {
			logger.Log.Error("Error accepting connection", "stack", err)
			break
		}
		clientAddr := newConn.RemoteAddr().String()
		// prepare the new QUIC CONNECTION obj
		qc := &QuicConn{
			Conn:            newConn,
			UpstreamServer:  qsi,
			ActiveStream:    make(ACTIVE_STREAM),
			IsAuthenticated: false,
		}

		qsi.ActiveConn[clientAddr] = qc

		qsi.Wg.Add(1)
		go qc.serve(baseCtx)

	}
	// terminate all the go routines to prevent information loss
	qsi.Wg.Wait()

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
