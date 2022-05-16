package pubsub

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SyncEventMsg struct {
	Address   common.Address `json:"address"`
	Block     uint64         `json:"block"`
	Timestamp uint64         `json:"timestamp"`
	Reserve0  *big.Int       `json:"reserve0"`
	Reserve1  *big.Int       `json:"reserve1"`
}

type PriceAlertMsg struct {
	ServerId    string  `json:"serverId"`
	PoolAddress string  `json:"poolAddress"`
	ChannelId   string  `json:"channelId"`
	PairName    string  `json:"pairName"`
	Change      float64 `json:"change"`
	Since       int64   `json:"since"`
}

type PriceSummaryMsg struct {
	ServerId    string  `json:"serverId"`
	PoolAddress string  `json:"poolAddress"`
	ChannelId   string  `json:"channelId"`
	PairName    string  `json:"pairName"`
	Date        string  `json:"date"`
	Open        float64 `json:"open"`
	Close       float64 `json:"close"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Change      float64 `json:"change"`
}
