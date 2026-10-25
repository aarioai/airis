package atype

import (
	"fmt"
	"github.com/aarioai/airis/pkg/types"
	"math"
	"strings"
)

// Exchange rates usually keep 4 decimal places, so amounts here are value * 10000.
// The database money() function supports signed money.
// SmallMoney: UNSIGNED INT range, about 420,000 yuan.
// Money:      BIGINT range, about ±10 billion yuan.

const (
	// 1 yuan = 100 fen = 1000 hao = 10000 money
	// 10 lu-coins = 100 lu-fen = 1000 lu-li = 1000 lu-hao = 10000 coin

	MoneyScale        uint8   = 4
	Cent              Money   = 100 // 0.01
	Cents                     = Cent
	Dime              Money   = 1000 // 0.1
	Dimes                     = Dime
	MoneyUnits        Money   = 10000 // any money unit, e.g., CNY 1, USD 1
	moneyUnitsFloat64 float64 = 10000.0

	Yuan    = MoneyUnits
	KYuan   = 1000 * Yuan     // CNY 1000
	WanYuan = 10000 * Yuan    // CNY 10,000
	MYuan   = 100 * WanYuan   // CNY 1000,000
	YiiYuan = 10000 * WanYuan // CNY 100,000,000

	Dollar   = MoneyUnits
	Dollars  = Dollar
	KDollar  = 1000 * Dollar // 1000 dollar
	KDollars = KDollar
	MDollar  = 1000 * KDollar // million dollar
	MDollars = MDollar
	BDollar  = 1000 * MDollar // billion dollar
	BDollars = BDollar

	MaxMoneyU32 Money = 1<<32 - 1
	MinMoney    Money = -1 << 63
	MaxMoney    Money = 1<<63 - 1
)

// For an int64, multiply directly, e.g. 100 * atype.Yuan
func MoneyUnitsX(n float64) Money { return Money(math.Round(n * moneyUnitsFloat64)) }
func YuanX(n float64) Money       { return MoneyUnitsX(n) }
func DollarX(n float64) Money     { return MoneyUnitsX(n) }

func (a Money) Int64() int64   { return int64(a) }
func (a Money) String() string { return types.FormatInt(a) }

// Compute the total price
func (a Money) MulN(n uint) Money { return a * Money(n) }

// Compute the discounted price
func (a Money) MulF(p float64, intHandler func(float64) float64) Money {
	return Money(intHandler(float64(a) * p))
}
func (a Money) MulRoundF(p float64) Money { return a.MulF(p, math.Round) }
func (a Money) MulCeilF(p float64) Money  { return a.MulF(p, math.Ceil) }
func (a Money) MulFloorF(p float64) Money { return a.MulF(p, math.Floor) }

// Compute the discounted price
func (a Money) Mul(p Decimal, intHandler func(float64) float64) Money {
	return Money(intHandler(float64(a*Money(p)) / decimalUnitsFloat64))
}
func (a Money) MulRound(p Decimal) Money { return a.Mul(p, math.Round) }
func (a Money) MulCeil(p Decimal) Money  { return a.Mul(p, math.Ceil) }
func (a Money) MulFloor(p Decimal) Money { return a * Money(p) / Money(decimalUnitsInt64) }

// Compute the pre-discount price
func (a Money) Div(p Decimal, intHandler func(float64) float64) Money {
	return Money(intHandler(float64(a.Int64()*decimalUnitsInt64) / float64(p)))
}
func (a Money) DivRound(p Decimal) Money { return a.Div(p, math.Round) }
func (a Money) DivCeil(p Decimal) Money  { return a.Div(p, math.Ceil) }
func (a Money) DivFloor(p Decimal) Money { return a * Money(decimalUnitsInt64) / Money(p) }

// Compute the unit price
func (a Money) DivN(n uint, intHandler func(float64) float64) Money {
	return a.Div(Decimal(n)*DecimalUnits, intHandler)
}
func (a Money) DivRoundN(n uint) Money { return a.DivN(n, math.Round) }
func (a Money) DivCeilN(n uint) Money  { return a.DivN(n, math.Ceil) }
func (a Money) DivFloorN(n uint) Money { return a.DivFloor(Decimal(n) * DecimalUnits) }

// Ratio of one charge to the total charge
func (a Money) Of(b Money, intHandler func(float64) float64) Decimal {
	return Decimal(intHandler(float64(a.Int64()*decimalUnitsInt64) / float64(b)))
}
func (a Money) OfRound(b Money) Decimal { return a.Of(b, math.Round) }
func (a Money) OfCeil(b Money) Decimal  { return a.Of(b, math.Ceil) }
func (a Money) OfFloor(b Money) Decimal { return Decimal(a) * DecimalUnits / Decimal(b) }

// Sign
func (a Money) Sign() string {
	if a > 0 {
		return "-"
	}
	return ""
}

// Precision
func (a Money) Precision() int {
	n := len(types.FormatInt(a.Int64()))
	if a < 0 {
		n--
	}
	return n
}

// Real-number form
func (a Money) ToReal() float64 { return float64(a) / moneyUnitsFloat64 }

// Integer part
func (a Money) Whole() int64  { return int64(a) / int64(Yuan) }
func (a Money) ToCent() int64 { return int64(a) / int64(Cent) }

// Fractional part. Represent it as a fraction, not an integer, because 123.0001 ---> 0001
func (a Money) Mantissa(withSign bool) float64 {
	mantissa := float64(a%MoneyUnits) / moneyUnitsFloat64
	if !withSign && mantissa < 0 {
		mantissa = -mantissa
	}
	return mantissa
}

// Style such as 1,000,000
func (a Money) FormatWhole(style *DecimalFormat) string {
	s := types.FormatInt(a.Whole())
	if style == nil {
		return s
	}
	return formatWhole(s, style.SegmentSize, style.Separator)
}

// @warn If rounding carries into the integer, keep only .999...; negatives are rounded on the absolute value
func (a Money) FormatMantissa(style *DecimalFormat) string {
	style = NewDecimalFormat(style)
	scale := int(style.Scale)

	s := types.FormatInt(a.Int64())
	g := len(s) - int(MoneyScale)
	if g > 0 {
		s = s[g:] // take the fractional part
	}
	var ok bool
	if s, ok = mantissaOk(s, scale, style.TrimScale); !ok {
		return s
	}
	if len(s) > scale {
		b := s[len(s)-1]
		if scale > 0 {
			b = s[scale-1]
		}
		k := style.ScaleRound.IsCeil() || (style.ScaleRound.IsRound() && strings.IndexByte("0123456789", b) > 4)
		s = s[:scale]
		if k {
			// 0.999... does not need to carry
			repeat9 := true
			for _, ss := range s {
				if ss != '9' {
					repeat9 = false
					break
				}
			}
			if !repeat9 {
				x, _ := types.ParseUint16("1" + s)
				s = types.FormatUint(x + 1)
			}
		}
		// s has changed
		if s, ok = mantissaOk(s, scale, style.TrimScale); !ok {
			return s
		}
	}
	if !style.TrimScale {
		s = padRight(s, "0", scale)
	}
	return "." + s
}

func (a Money) Format(style *DecimalFormat) string {
	style = NewDecimalFormat(style)
	fmt.Println(a.FormatWhole(style), a.FormatMantissa(style))
	return a.FormatWhole(style) + a.FormatMantissa(style)
}

// No parameters; provided for Go templates
func (a Money) FmtWhole() string {
	return a.FormatWhole(nil)
}
func (a Money) FmtMantissa() string {
	return a.FormatMantissa(&DecimalFormat{Scale: 2})
}
func (a Money) Fmt() string {
	return a.Format(&DecimalFormat{Scale: 2})
}
