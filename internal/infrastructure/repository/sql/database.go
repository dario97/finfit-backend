package sql

import "gorm.io/gorm"

type Database interface {
	First(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}
