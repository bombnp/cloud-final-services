package models

type Token struct {
	Address string `gorm:"primarykey;column:address"`
	Symbol  string `gorm:"column:symbol"`
	Icon    string `gorm:"column:icon"`
	Name    string `gorm:"column:name"`
}
