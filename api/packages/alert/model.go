package alert

import "time"

type Logger struct {
	Message string `json:"message"`
}

type AlertSummary struct {
	High     float64
	HighTime time.Time
	Low      float64
	LowTime  time.Time
	Change   float64
}

type AlertResponse struct {
	Address string  `json:"address"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Change  float64 `json:"change"`
}
