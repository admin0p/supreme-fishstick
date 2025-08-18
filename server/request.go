package server

type REQUEST_OBJ struct {
	From           string
	To             string
	Type           string
	MessageFormat  string
	MessageId      string
	Payload        *[]byte
	UpstreamServer *QUIC_SERVER_INSTANCE
}
