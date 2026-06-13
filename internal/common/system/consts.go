package system

const (
	ResponseKey = "response"
)

type Code int

const (
	SuccessCode Code = 20000000
	ClientError Code = 40000000
	ServerError Code = 50000000
)
