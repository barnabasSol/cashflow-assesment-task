package paymentmaking

import "time"

type CreatePaymentRequest struct {
	Ref      string `json:"ref"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreatePaymentResponse struct {
	PaymentId string `json:"payment_id"`
	Status    int    `json:"status"`
}

type GetPaymentStatusResponse struct {
	Amount    string    `json:"amount"`
	Currency  string    `json:"currency"`
	Ref       string    `json:"ref"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
