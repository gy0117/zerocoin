package engine

import "github.com/pkg/errors"

var (
	TradePairError   = errors.New("tradePair error")
	TimeoutError     = errors.New("timeout error")
	OrderSideError   = errors.New("order side error")
	OrderTypeError   = errors.New("order type error")
	CancelOrderError = errors.New("cancel order error")
)
