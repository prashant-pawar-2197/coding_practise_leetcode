package main

import (
	"fmt"
	"net"
	"strconv"
)

func hexToIP(hex string) (string, error) {
	// Convert the hexadecimal string to an integer
	num, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return "", err
	}

	// Convert the integer to an IP address
	ip := net.IPv4(byte(num>>24), byte(num>>16), byte(num>>8), byte(num))
	return ip.String(), nil
}

func main() {
	hex := "0101A8C0" // Example hexadecimal number
	ip, err := hexToIP(hex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("IP Address:", ip)
}