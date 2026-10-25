package atype

import (
	"github.com/aarioai/airis/pkg/types"
	"math"
)

// virtual money

const (
	VCent              VMoney  = 100  // 1 lu-fen
	VDime              VMoney  = 1000 // 1 lu-coin (jiao)
	VMoneyUnits        VMoney  = 10000
	vmoneyUnitsFloat64 float64 = 10000.0
	vmoneyUnitsInt64   int64   = 10000

	MaxVMoneyU32 VMoney = 1<<32 - 1
	MinVMoney    VMoney = -1 << 63
	MaxVMoney    VMoney = 1<<63 - 1
)

func VmoneyUnitsX(n float64) VMoney { return VMoney(math.Round(n * vmoneyUnitsFloat64)) }

func (c VMoney) Int64() int64   { return int64(c) }
func (c VMoney) String() string { return types.FormatInt(c) }

// Real-number form
func (c VMoney) Real() float64 { return float64(c) / vmoneyUnitsFloat64 }

// Gold coins offset against goods
//
//	.Off()  ==  .Offset(100*Percent)
func (c VMoney) Offset(rate Decimal) Money { return Money(c).MulFloor(rate) }
func (c VMoney) Off() Money                { return Money(c) }

// For use in calculations
func (c VMoney) P() Money { return Money(c) }

// Purchase quantity of goods corresponding to this amount of gold coins.
// In a gold-coin SKU, exchange_rate<Decimal> is the number of gold coins per item.
// Gold coins redeemable per item: unit_coin<VMoney> = exchange_rate.ToReal * VMoneyUnits = exchange_rate * VMoneyUnits / DecimalUnits
// Total redeemed amount: total_coin = unit_coin * qty
// qty = total_coin / unit_coin = (total_coin * DecimalUnits) / (exchange_rate * VMoneyUnits)
func ExchangeVMoney(exchangeRate Decimal, qty uint) VMoney {
	return VMoney(exchangeRate.Int64() * int64(qty) * vmoneyUnitsInt64 / decimalUnitsInt64)
}
func (c VMoney) ExchangeQty(exchangeRate Decimal, intHandler func(float64) float64) uint {
	return uint(intHandler(float64(c.Int64()*decimalUnitsInt64) / float64(exchangeRate.Int64()*vmoneyUnitsInt64)))
}
func (c VMoney) ExchangeQtyCeil(exchangeRate Decimal) uint {
	return c.ExchangeQty(exchangeRate, math.Ceil)
}
func (c VMoney) ExchangeQtyFloor(exchangeRate Decimal) uint {
	return c.ExchangeQty(exchangeRate, math.Floor)
}
func (a Money) VMoney() VMoney { return VMoney(a) }
