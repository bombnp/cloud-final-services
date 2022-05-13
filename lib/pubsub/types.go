package pubsub

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SyncEventMsg struct {
	Address  common.Address `json:"address"`
	Block    uint64         `json:"block"`
	Reserve0 *big.Int       `json:"reserve0"`
	Reserve1 *big.Int       `json:"reserve1"`
}
