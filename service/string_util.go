package service

import (
	"fmt"
	"strconv"
	"unicode"
)

// IncrementStringEnd appends "1" to the input string s, or increments the number at the end if present.
func IncrementStringEnd(s string) string {
	length := len(s)
	if length == 0 {
		return "1"
	}

	index := length
	for index > 0 && unicode.IsDigit(rune(s[index-1])) {
		index--
	}

	if index == length {
		return s + " 1"
	}

	numberPart := s[index:]
	prefix := s[:index]
	number, err := strconv.Atoi(numberPart)
	if err != nil {
		// This should not happen since we have already checked that the suffix is numeric
		fmt.Println("Error converting string to integer:", err)
		return s + " 1"
	}

	// Increment the number and concatenate it back with the prefix
	return fmt.Sprintf("%s%d", prefix, number+1)
}
