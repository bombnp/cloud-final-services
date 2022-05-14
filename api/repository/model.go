package repository

type Token struct {
	Address string `gorm:"column:address"`
	Symbol  string `gorm:"column:symbol"`
	Icon    string `gorm:"column:icon"`
}

type Pair struct {
	PoolAddress  string `gorm:"column:pool_address"`
	BaseAddress  string `gorm:"column:base_address"`
	QuoteAddress string `gorm:"column:quote_address"`
	IsBaseToken0 bool   `gorm:"column:is_base_token0"`
}
