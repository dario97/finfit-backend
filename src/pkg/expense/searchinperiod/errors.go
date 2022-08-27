package searchinperiod

type invalidArgumentsError struct {
	msg string
}

func (receiver invalidArgumentsError) Error() string {
	return receiver.msg
}
