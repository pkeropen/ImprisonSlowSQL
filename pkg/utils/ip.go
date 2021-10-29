package utils

import (
	"strconv"
	"strings"
)

func ValidIPAddress(IP string) string {
	if CheckIPv4(IP) {
		return "IPv4"
	}
	if CheckIPv6(IP) {
		return "IPv6"
	}
	return "Neither"
}


func CheckIPv4(IP string) bool {
	// 字符串这样切割
	strs := strings.Split(IP, ".")
	if len(strs) != 4 {
		return false
	}
	for _, s := range strs {
		if len(s) == 0 || (len(s) > 1 && s[0] == '0') {
			return false
		}
		// 直接访问字符串的值
		if s[0] < '0' || s[0] > '9' {
			return false
		}
		// 字符串转数字
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if n < 0 || n > 255 {
			return false
		}
	}
	return true
}

func CheckIPv6(IP string) bool {
	strs := strings.Split(IP, ":")
	if len(strs) != 8 {
		return false
	}
	for _, s := range strs {
		if len(s) <= 0 || len(s) > 4 {
			return false
		}
		for i := 0; i < len(s); i++ {
			if s[i] >= '0' && s[i] <= '9' {
				continue
			}
			if s[i] >= 'A' && s[i] <= 'F' {
				continue
			}
			if s[i] >= 'a' && s[i] <= 'f' {
				continue
			}
			return false
		}
	}
	return true
}
