package variables

var (
	ERROR_REDIS         []byte = []byte("-Error message\r\n")
	ERROR_TIMEOUT_REDIS []byte = []byte("-Timeout\r\n")
	ERROR_READ_REDIS    []byte = []byte("-Error while reading\r\n")
	REDIS_END           []byte = []byte("\r\n")
	REDIS_QUIT          []byte = []byte("*1\r\n$4\r\nquit\r\n")

	REDIS_ONE_COMMAND string = "*1\r\n$%d\r\n%s\r\n"
	REDIS_TWO_COMMAND string = "*2\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n"

	SENTINEL_COMMAND string = "sentinel get-master-addr-by-name %s\n"

	BUILDING_CODE         string = "BUILDING_CODE"
	UTILITY_BUILDING_CODE string = "UTILITY_BUILDING_CODE"
	PROCEDURE_CODE        string = "PROCEDURE_CODE"
	DOCUMENT_CODE         string = "DOCUMENT_CODE"

	SIMPLE_STRING byte = '+'
	BULK_STRING   byte = '$'
	INTEGER       byte = ':'
	ARRAY         byte = '*'
	ERROR         byte = '-'

	BUFFER_SIZE int = 256 * 1024
)
