package pair

type Logger struct {
	Message string `json:"message"`
}

type GetPairsResponse struct {
	PoolAddress  string `json:"pool_address"`
	PoolName     string `json:"pool_name"`
	IsBaseToken0 bool   `json:"is_base_token0"`
}
