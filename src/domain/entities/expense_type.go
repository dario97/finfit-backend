package entities

type ExpenseType struct {
	id   int64
	name string
}

func NewExpenseType(name string) ExpenseType {
	return ExpenseType{
		name: name,
	}
}

func NewExpenseTypeWithId(id int64, name string) ExpenseType {
	return ExpenseType{
		id:   id,
		name: name,
	}
}

func (et ExpenseType) Id() int64 {
	return et.id
}

func (et ExpenseType) Name() string {
	return et.name
}
