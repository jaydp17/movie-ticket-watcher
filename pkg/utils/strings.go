package utils

import (
	"fmt"
	"math/rand"
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

//RandomString - Generate a random string of A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}
