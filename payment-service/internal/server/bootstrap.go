package server

import paymentmaking "payment-service/internal/features/payment-making"

func (s *Server) bootstrap() error {
	pmr := paymentmaking.NewRepository(s.db)
	pms := paymentmaking.NewService(pmr, s.rmq)
	_ = pms
	return nil
}
