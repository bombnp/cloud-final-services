package subscribe

type AlertSubscribeRequest struct {
	ServerId    string `json:"server_id"`
	PoolAddress string `json:"pool"`
	ChannelId   string `json:"channel_id"`
}

type Logger struct {
	Message string `json:"message"`
}
