package ethutils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type LiquidityPair interface {
	SwapEventID() common.Hash
	SyncEventID() common.Hash

	UnmarshalSwapEvent(data []byte) (*SwapEvent, error)
	UnmarshalSyncEvent(data []byte) (*SyncEvent, error)
}

const DexTypeUniswapV2 = "uniswapv2"

var lps = map[string]LiquidityPair{
	DexTypeUniswapV2: newUniswapV2LP(),
}

func GetLp(dexType string) (LiquidityPair, error) {
	lp, ok := lps[dexType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[GetLp]: unknown dex type: %s", dexType))
	}
	return lp, nil
}

var swapTopics = []common.Hash{lps[DexTypeUniswapV2].SwapEventID()}
var syncTopics = []common.Hash{lps[DexTypeUniswapV2].SyncEventID()}

var swapAndSyncTopics = append(swapTopics, syncTopics...)

func GetSwapTopics() []common.Hash {
	return swapTopics
}

func GetSyncTopics() []common.Hash {
	return syncTopics
}

func GetSwapAndSyncTopics() []common.Hash {
	return swapAndSyncTopics
}
