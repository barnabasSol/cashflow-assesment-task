package paymentmaking

type CreatePaymentRequest struct {
	Ref      string `json:"ref"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type CreatePaymentResponse struct {
	PaymentId string `json:"payment_id"`
	Status    int    `json:"status"`
}
