package server

type ProtoHandler struct {
	Db string // Database connection string or identifier
}

// this is not needed...will be a internal function that will be used by the server to process the stream
// func (ph *ProtoHandler) ProcessStream(ctx context.Context, sfStreamHandler *SF_STREAM_HANDLER) error {
// 	ackPackage := &dataframe.ACK_FRAME{
// 		StreamId:  int32(sfStreamHandler.Stream.StreamID()),
// 		PackId:    1,
// 		AckStatus: true,
// 	}
// 	//receive package
// 	result := dataframe.MESSAGE_FRAME{}
// 	_, err := sfStreamHandler.PackageHandler.ReceivePackage(ctx, sfStreamHandler.Stream, &result)
// 	if err != nil {
// 		fmt.Println("Failed to receive package:", err)
// 		return err
// 	}

// 	// test print
// 	fmt.Println("Received package:", result.ProtoReflect().Type().Descriptor().Name())

// 	_, err = sfStreamHandler.PackageHandler.SendPackage(ctx, sfStreamHandler.Stream, ackPackage)
// 	if err != nil {
// 		fmt.Println("Failed to send ACK package:", err)
// 		return err
// 	}

// 	return nil
// }
