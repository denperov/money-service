package models

import (
	"testing"
)

func TestMoney_String(t *testing.T) {
	t.Run("0.00", func(t *testing.T) {
		money0 := Money{0}
		if money0.String() != "0.00" {
			t.Fail()
		}
	})
	t.Run("0.01", func(t *testing.T) {
		money1 := Money{1}
		if money1.String() != "0.01" {
			t.Fail()
		}
	})
	t.Run("0.10", func(t *testing.T) {
		money10 := Money{10}
		if money10.String() != "0.10" {
			t.Fail()
		}
	})
	t.Run("1.00", func(t *testing.T) {
		money100 := Money{100}
		if money100.String() != "1.00" {
			t.Fail()
		}
	})
	t.Run("max", func(t *testing.T) {
		moneyMax := MaxMoney()
		if moneyMax.String() != "45035996273704.95" {
			t.Fail()
		}
	})
}

func TestMoneyFromString(t *testing.T) {
	t.Run("0.00", func(t *testing.T) {
		x, err := MoneyFromString("0.00")
		if err != nil {
			t.Fail()
		}
		if x != (Money{0}) {
			t.Fail()
		}
	})
	t.Run("0.01", func(t *testing.T) {
		x, err := MoneyFromString("0.01")
		if err != nil {
			t.Fail()
		}
		if x != (Money{1}) {
			t.Fail()
		}
	})
	t.Run("0.10", func(t *testing.T) {
		x, err := MoneyFromString("0.10")
		if err != nil {
			t.Fail()
		}
		if x != (Money{10}) {
			t.Fail()
		}
	})
	t.Run("1.00", func(t *testing.T) {
		x, err := MoneyFromString("1.00")
		if err != nil {
			t.Fail()
		}
		if x != (Money{100}) {
			t.Fail()
		}
	})
	t.Run("max", func(t *testing.T) {
		x, err := MoneyFromString("45035996273704.95")
		if err != nil {
			t.Fail()
		}
		if x != MaxMoney() {
			t.Fail()
		}
	})
}

func TestMoneySum(t *testing.T) {
	t.Run("0+0", func(t *testing.T) {
		x, err := MoneySum(ZeroMoney(), ZeroMoney())
		if err != nil {
			t.Fail()
		}
		if x != ZeroMoney() {
			t.Fail()
		}
	})
	t.Run("1+1", func(t *testing.T) {
		x, err := MoneySum(Money{1}, Money{1})
		if err != nil {
			t.Fail()
		}
		if x != (Money{2}) {
			t.Fail()
		}
	})
	t.Run("max+0", func(t *testing.T) {
		x, err := MoneySum(MaxMoney(), ZeroMoney())
		if err != nil {
			t.Fail()
		}
		if x != MaxMoney() {
			t.Fail()
		}
	})
	t.Run("max+1", func(t *testing.T) {
		_, err := MoneySum(MaxMoney(), Money{1})
		if err == nil {
			t.Fail()
		}
	})
}

func TestMoneyDiff(t *testing.T) {
	t.Run("0-0", func(t *testing.T) {
		x, err := MoneyDiff(ZeroMoney(), ZeroMoney())
		if err != nil {
			t.Fail()
		}
		if x != ZeroMoney() {
			t.Fail()
		}
	})
	t.Run("1+1", func(t *testing.T) {
		x, err := MoneyDiff(Money{1}, Money{1})
		if err != nil {
			t.Fail()
		}
		if x != ZeroMoney() {
			t.Fail()
		}
	})
	t.Run("min-0", func(t *testing.T) {
		x, err := MoneyDiff(MinMoney(), ZeroMoney())
		if err != nil {
			t.Fail()
		}
		if x != MinMoney() {
			t.Fail()
		}
	})
	t.Run("min-1", func(t *testing.T) {
		_, err := MoneyDiff(MinMoney(), Money{1})
		if err == nil {
			t.Fail()
		}
	})
}
