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
	Address string
	High    float64
	Low     float64
	Change  float64
}
