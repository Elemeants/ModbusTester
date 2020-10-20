package utils

import "fmt"

func PrintBuffer(buf []byte) {
	for _, v := range buf {
		fmt.Printf(" 0x%02X", v)
	}
	fmt.Println()
}
