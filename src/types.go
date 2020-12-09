package main

type (
	// WebSocket的信息流数据类型
	stream struct {
		Event   string      `json:"event"`
		Payload interface{} `json:"payload"`
	}

	Tag struct {
		name  string
		count int
	}
)
