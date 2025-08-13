package main

import (
	dataframe "github.com/admin0p/supreme-fishstick/proto"
	"google.golang.org/protobuf/proto"
)

// func main() {

// 	quicServer := &server.QUIC_SERVER_INSTANCE{
// 		HostName: "localhost",
// 		Port:     4242,
// 	}

// 	quicServer.StartServer(nil, 1)
// }

type PROTO_MESSAGE_DUMMY struct {
	proto.Message
}

func main() {

	msg := dataframe.MESSAGE_FRAME{
		From:    "client1",
		Message: "Hello, World!",
	}

	x, err := proto.Marshal(&msg)
	if err != nil {
		panic(err)
	}
	result := PROTO_MESSAGE_DUMMY{}
	err = proto.Unmarshal(x, &result)
	if err != nil {
		panic(err)
	}
	println("Message:", result.Message.ProtoReflect().Type())
}
