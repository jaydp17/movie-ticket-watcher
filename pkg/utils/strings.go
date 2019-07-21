package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func ToFloat(str string) float64 {
	floatNum, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Printf("error while parsing to float: %v\nError: %v\n", str, err)
		return 0.0
	}
	return floatNum
}

func KeepJustAlphaNumeric(str string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return reg.ReplaceAllString(str, "")
}
