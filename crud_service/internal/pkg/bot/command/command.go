package command

import "context"

var (
	MsgInvalidArguments = "invalid arguments"
	MsgValidationError  = "validation error"
	MsgInternalError    = "internal error"
)

type Interface interface {
	Name() string
	Description() string
	Process(ctx context.Context, args string) string
}
