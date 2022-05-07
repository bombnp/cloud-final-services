package services

type SubscribeRequest struct {
	ServerId    string `json:"server_id"`
	PoolAddress string `json:"pool"`
}

type Logger struct {
	Message string `json:"message"`
}
