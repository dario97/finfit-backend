package custom_errors

type InternalError struct {
	Msg string
}

func (receiver InternalError) Error() string {
	return receiver.Msg
}
