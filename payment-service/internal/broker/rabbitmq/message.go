package broker

import "encoding/json"

type Message struct {
	Exchange string
	Topic    string          `json:"topic"`
	Data     json.RawMessage `json:"data"`
}

type NewPayment struct {
	PaymentId string `json:"id"`
}
