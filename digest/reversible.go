package digest

import "encoding/base64"

func Base64(plainText []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(plainText))
}
func Debase64(cipherText []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(cipherText))
}
