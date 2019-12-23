package models

import (
	"errors"
	"fmt"
	"strconv"
)

type Money struct{ int64 }

func ZeroMoney() Money { return Money{0} }
func MaxMoney() Money  { return Money{4503599627370495} }  // 52-bit (2^52-1)
func MinMoney() Money  { return Money{-4503599627370495} } // 52-bit (2^52-1) * -1

func MoneyFromString(val string) (Money, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return Money{}, err
	}
	if f > float64(MaxMoney().int64) {
		return Money{}, errors.New("string cannot be represented as money type")
	}
	return Money{int64(f * 100)}, nil
}

func (m Money) String() string {
	return fmt.Sprintf("%d.%02d", m.int64/100, m.int64%100)
}

func (m *Money) UnmarshalText(text []byte) (err error) {
	*m, err = MoneyFromString(string(text))
	return
}

func (m Money) MarshalText() (text []byte, err error) {
	return []byte(m.String()), nil
}

func MoneyIsPositive(a Money) bool {
	return a.int64 > 0
}

func MoneyLess(a, b Money) bool {
	return a.int64 < b.int64
}

func MoneyDiff(a, b Money) (Money, error) {
	diff := a.int64 - b.int64
	if diff < MinMoney().int64 {
		return Money{}, errors.New("money diff overflow")
	}
	return Money{diff}, nil
}

func MoneySum(a, b Money) (Money, error) {
	sum := a.int64 + b.int64
	if sum > MaxMoney().int64 {
		return Money{}, errors.New("money sum overflow")
	}
	return Money{sum}, nil
}
