package models

type Pair struct {
	PoolAddress  string `gorm:"primarykey;column:pool_address"`
	BaseAddress  string `gorm:"column:base_address"`
	QuoteAddress string `gorm:"column:quote_address"`
	IsBaseToken0 bool   `gorm:"column:is_base_token0"`
}
