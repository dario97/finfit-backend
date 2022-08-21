package add

type UnexpectedError struct {
	Msg string
}

func (receiver UnexpectedError) Error() string {
	return receiver.Msg
}

type InvalidExpenseTypeError struct {
	Msg string
}

func (receiver InvalidExpenseTypeError) Error() string {
	return receiver.Msg
}
