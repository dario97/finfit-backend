package custom_errors

type InvalidExpenseTypeError struct {
	Msg string
}

func (receiver InvalidExpenseTypeError) Error() string {
	return receiver.Msg
}
