package ethutils

import "math/big"

type SwapEvent struct {
	// positive means out and negative means in
	Amount0 *big.Int
	// positive means out and negative means in
	Amount1 *big.Int
}

type rawSwapEvent struct {
	Amount0In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	Amount1In  *big.Int
}

type SyncEvent struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
}
