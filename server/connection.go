package server

import (
	"context"

	"github.com/admin0p/supreme-fishstick/logger"
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

type QuicConn struct {
	Conn            *quic.Conn
	UpstreamServer  *QUIC_SERVER_INSTANCE
	ActiveStream    ACTIVE_STREAM
	IsAuthenticated bool
}

func (qc *QuicConn) serve(ctx context.Context) {
	// authenticate the client 1st
	authStream, err := qc.Conn.OpenStreamSync(ctx)
	if err != nil {
		logger.Log.Error("Failed to open auth stream", "stack", err)
		return
	}
	err = qc.authenticate(ctx, authStream)
	if err != nil {
		return
	}

	qc.ActiveStream["SESSION"], err = qc.Conn.AcceptStream(ctx)
	if err != nil {
		logger.Log.Error("Failed to accept session stream", "stack", err)
		return
	}

	for {
		qc.RequestHandler(ctx)
	}

}

func (qc *QuicConn) authenticate(ctx context.Context, authStream *quic.Stream) error {
	// read the buffer that constructs the request obj for auth

	qc.ActiveStream["AUTH"] = authStream
	defer authStream.Close()

	rawAuthPayload, err := readBuffer(authStream)
	if err != nil {
		return err
	}
	authFrame := &dataframe.CLIENT_AUTH_REQUEST_FRAME{}
	err = proto.Unmarshal(*rawAuthPayload, authFrame)
	if err != nil {
		logger.Log.Error("Failed to unmarshal auth frame", "stack", err)
		return err
	}

	// prepare the request obj
	newAuthReq := SF_REQ_CONTEXT{
		Ctx:            ctx,
		ProtoPkg:       authFrame,
		Payload:        rawAuthPayload,
		upstreamServer: qc.UpstreamServer,
	}

	// call the auth handler
	qc.IsAuthenticated, err = qc.UpstreamServer.Handler.AuthHandler(&newAuthReq)
	if err != nil || !qc.IsAuthenticated {
		logger.Log.Info("Authentication failed", "stack", err)
		return err
	}

	return nil
}

func (qc *QuicConn) RequestHandler(ctx context.Context) error {
	// read form the buffer
	rawPayload, err := readBuffer(qc.ActiveStream["SESSION"])
	if err != nil {
		return err
	}
	messageFrame := &dataframe.MESSAGE_FRAME{}
	err = proto.Unmarshal(*rawPayload, messageFrame)
	if err != nil {
		logger.Log.Error("Failed to unmarshal message frame", "stack", err)
		return err
	}
	// construct message Request Obj
	newSessionReqObj := &SF_REQ_CONTEXT{
		Ctx:            ctx,
		ProtoPkg:       messageFrame,
		Payload:        rawPayload,
		upstreamServer: qc.UpstreamServer,
	}
	// call the message handler
	qc.UpstreamServer.Handler.messageHandler(newSessionReqObj)
	return nil
}
