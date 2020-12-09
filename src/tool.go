package main

import (
	"unicode"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// 监测数组是否包含某自负
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// 监测字符串中是否包含中文
func hasHan(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}

	return false
}
