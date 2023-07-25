package digest

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(plainText string) string {
	md5er := md5.New()
	md5er.Write([]byte(plainText))
	return hex.EncodeToString(md5er.Sum(nil))
}
