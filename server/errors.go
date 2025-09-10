package server

import "github.com/admin0p/supreme-fishstick/logger"

type ERROR_CODE uint64

/*

 pre defied errors

 CONNECTION ERROR
 1. Server bind error
 2. connection accept error
 3. Stream accept error
 4. stream start error
 5. connection termination error
*/

const (
	ADDR_BIND            ERROR_CODE = 1
	CONN_DECLINED        ERROR_CODE = 2
	STREAM_DECLINED      ERROR_CODE = 3
	STREAM_START_ERR     ERROR_CODE = 4
	CONN_TERMINATION_ERR ERROR_CODE = 5
)

/*
	 BUFFER ERROR
	 1. buffer read error
		-> Max size exceeded error
	 2. buffer send error
		-> Max size exceeded error
*/
const (
	BUFF_READ_FAILED   ERROR_CODE = 6
	BUFF_SIZE_EXCEEDED ERROR_CODE = 7
	BUFF_WRITE_FAILED  ERROR_CODE = 8
)

/*
REQUEST ERROR
1. Stream does not exist
2. Connection does not exists
3. Invalid stream name

*/

const (
	INVALID_STREAM_REF  ERROR_CODE = 9
	INVALID_CONN        ERROR_CODE = 10
	INVALID_STREAM_NAME ERROR_CODE = 11
)

type SF_ERROR struct {
	Err      error
	ErrCode  ERROR_CODE
	BaseType string
	Desc     string
}

func (se *SF_ERROR) Error() string {
	return se.Err.Error()
}

func NewErr(err error, code ERROR_CODE, desc string) *SF_ERROR {
	return &SF_ERROR{
		Err:      err,
		ErrCode:  code,
		Desc:     desc,
		BaseType: "None",
	}

}

func (se *SF_ERROR) LogError(msg string) {
	if msg == "" {
		logger.Log.Error(se.Desc, "ERR_CODE", se.ErrCode, "stack", se.Err)
		return
	}
	logger.Log.Error(msg, "ERR_CODE", se.ErrCode, "stack", se.Err)
}
