package sql

import "gorm.io/gorm"

// TODO: Esta interfaz es algo inutil por el momento, los metodos no deberian devolver el DB de gorm, deberian devolver esta interfaz (Solucionar)
type Database interface {
	First(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Table(name string, args ...interface{}) (tx *gorm.DB)
}
