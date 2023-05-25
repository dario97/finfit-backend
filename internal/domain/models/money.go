package models

type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) *Money {
	return &Money{amount: amount, currency: currency}
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}
