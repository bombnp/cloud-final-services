package repository

type Token struct {
	Address string `gorm:"column:address"`
	Symbol  string `gorm:"column:symbol"`
	Icon    string `gorm:"column:icon"`
	Name    string `gorm:"column:name"`
}

type Pair struct {
	PoolAddress  string `gorm:"column:pool_address"`
	BaseAddress  string `gorm:"column:base_address"`
	QuoteAddress string `gorm:"column:quote_address"`
	IsBaseToken0 bool   `gorm:"column:is_base_token0"`
	Name         string `gorm:"column:name"`
}

type PairSubscription struct {
	Id          int    `gorm:"column:id"`
	ServerId    string `gorm:"column:server_id"`
	PoolAddress string `gorm:"column:pool_address"`
	Type        string `gorm:"column:type"`
	ChannelId   string `gorm:"column:channel_id"`
}
