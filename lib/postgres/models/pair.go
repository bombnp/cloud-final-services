package models

type Pair struct {
	PoolAddress      string  `gorm:"primarykey;column:pool_address"`
	BaseAddress      string  `gorm:"column:base_address"`
	QuoteAddress     string  `gorm:"column:quote_address"`
	IsBaseToken0     bool    `gorm:"column:is_base_token0"`
	Price            float64 `gorm:"column:price"`
	TwentyFourChange float64 `gorm:"column:24h_change"`
	TwentyFourVolume float64 `gorm:"column:24h_volume"`
	TwentyFourHigh   float64 `gorm:"column:24h_high"`
	TwentyFourLow    float64 `gorm:"column:24h_low"`
}
