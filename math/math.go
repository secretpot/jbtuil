package math

import (
	"fmt"
	"math"
	"strconv"
)

func Round(x float64, ndigits int) float64 {
	n := "%." + strconv.Itoa(ndigits) + "f"
	num, err := strconv.ParseFloat(fmt.Sprintf(n, x), 64)
	if err != nil {
		return x
	}
	return num
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Pow(math.E, -x))
}
