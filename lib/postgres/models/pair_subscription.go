package models

type PairSubscription struct {
	Id          int    `gorm:"primarykey;column:id;autoIncrement"`
	ServerId    string `gorm:"column:server_id"`
	PoolAddress string `gorm:"column:pool_address"`
	Type        string `gorm:"column:type"`
	ChannelId   string `gorm:"column:channel_id"`
}
