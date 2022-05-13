package repository

type Token struct {
	Address string `gorm:"address"`
	Symbol  string `gorm:"symbol"`
	Icon    string `gorm:"icon"`
}

type Pair struct {
	PoolAddress  string `gorm:"pool_address"`
	BaseAddress  string `gorm:"base_address"`
	QuoteAddress string `gorm:"quote_address"`
	IsBaseToken0 bool   `gorm:"is_base_token0"`
}
