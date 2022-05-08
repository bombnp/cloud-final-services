package repository

type Token struct {
	Address string `gorm:"address"`
	Symbol  string `gorm:"symbol"`
	Icon    string `gorm:"icon"`
}
