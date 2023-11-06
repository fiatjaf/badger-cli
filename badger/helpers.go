package badger

import (
	"encoding/hex"
)

func bytesToString(buf []byte) string {
	n := len(buf)
	readableCharacters := 0
	for _, b := range buf {
		if b >= 33 && b <= 126 {
			readableCharacters++
		}
	}
	if readableCharacters > n/2 {
		return string(buf)
	}
	return "0x" + hex.EncodeToString(buf)
}
