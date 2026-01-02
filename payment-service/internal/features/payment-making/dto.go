package paymentmaking

import "time"

type CreatePaymentRequest struct {
	Ref      string `json:"ref"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreatePaymentResponse struct {
	PaymentId string    `json:"ref"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
