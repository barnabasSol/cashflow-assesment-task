package paymentmaking

import (
	"errors"
	"strconv"
	"strings"
)

func parseAmountToMinorUnits(amountStr string) (int64, error) {
	amountStr = strings.TrimSpace(amountStr)
	if amountStr == "" {
		return 0, errors.New("amount is required")
	}

	parts := strings.Split(amountStr, ".")
	if len(parts) > 2 {
		return 0, errors.New("invalid amount format")
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, errors.New("invalid amount")
	}

	var cents int64 = 0

	if len(parts) == 2 {
		dec := parts[1]

		if len(dec) > 2 {
			return 0, errors.New("amount has more than 2 decimal places")
		}

		if len(dec) == 1 {
			dec += "0"
		}

		cents, err = strconv.ParseInt(dec, 10, 64)
		if err != nil {
			return 0, errors.New("invalid decimal amount")
		}
	}

	return whole*100 + cents, nil
}
