package services

type SubscribeRequest struct {
	ServerId    string `json:"server_id"`
	PoolAddress string `json:"pool"`
}

type Logger struct {
	Message string `json:"message"`
}

type PairResponse struct {
	PoolAddress  string `json:"pool_address"`
	PoolName     string `json:"pool_name"`
	IsBaseToken0 bool   `json:"is_base_token0"`
}
