package paymentmaking

import "errors"

var ValidCurrencies = map[string]struct{}{
	"USD": {},
	"ETB": {},
}

func (p CreatePaymentRequest) Validate() error {
	if len(p.Currency) != 3 {
		return errors.New("currency must be a 3-letter code")
	}
	if _, ok := ValidCurrencies[p.Currency]; !ok {
		return errors.New("not a valid currency")
	}
	if len(p.Ref) == 0 {
		return errors.New("reference must not be empty")
	}
	return nil
}
