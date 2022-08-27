package searchinperiod

type invalidArgumentsError struct {
	msg string
}

func (receiver invalidArgumentsError) Error() string {
	return receiver.msg
}

type unexpectedError struct {
	msg string
}

func (receiver unexpectedError) Error() string {
	return receiver.msg
}
