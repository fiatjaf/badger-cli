package badger

import (
	"encoding/hex"
	"strings"
)

func bytesToString(buf []byte) string {
	n := len(buf)
	readableCharacters := 0
	for _, b := range buf {
		if b >= 33 && b <= 126 {
			readableCharacters++
		}
	}
	if readableCharacters > n*2/3 {
		return string(buf)
	}
	return "0x" + hex.EncodeToString(buf)
}

func stringToBytes(str string) []byte {
	if strings.HasPrefix(str, "0x") {
		b, err := hex.DecodeString(str[2:])
		if err == nil {
			return b
		}
	}

	return []byte(str)
}
