package command

var (
	MsgInvalidArguments = "invalid arguments"
	MsgValidationError  = "validation error"
	MsgInternalError    = "internal error"
)

type Interface interface {
	Name() string
	Description() string
	Process(args string) string
}
