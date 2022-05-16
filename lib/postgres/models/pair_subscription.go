package models

type PairSubscription struct {
	Id          int    `json:"id" gorm:"primarykey;column:id;autoIncrement"`
	ServerId    string `json:"server_id" gorm:"column:server_id"`
	PoolAddress string `json:"pool_address" gorm:"column:pool_address"`
	Type        string `json:"type" gorm:"column:type"`
	ChannelId   string `json:"channel_id" gorm:"column:channel_id"`
}
