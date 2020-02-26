package variables

import "errors"

const (
	SENTINEL_COMMAND string = "sentinel get-master-addr-by-name %s\n"

	BUFFER_SIZE int = 1024 * 256

	SIMPLE_STRING byte = '+'
	INTEGER       byte = ':'
	ERROR         byte = '-'
	BULK_STRING   byte = '$'
	ARRAY         byte = '*'
)

var (
	REDIS_QUIT []byte = []byte("*1\r\n$4\r\nquit\r\n")

	ERROR_TIMEOUT_REDIS []byte = []byte("-Timeout\r\n")
	ERROR_READ_REDIS    []byte = []byte("-Error reading\r\n")

	ErrTimeOut       error = errors.New("Timeout")
	ErrWriting       error = errors.New("Writing error")
	ErrInvalidSyntax error = errors.New("Invalid syntax")
)
