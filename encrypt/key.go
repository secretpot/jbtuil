package encrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// GenerateRSAKey 生成RSA私钥和公钥
func GenerateRSAKey(bits int) ([]byte, []byte, error) {
	var (
		privateKeyContentBuffer bytes.Buffer = *bytes.NewBuffer(make([]byte, bits))
		publicKeyContentBuffer  bytes.Buffer = *bytes.NewBuffer(make([]byte, bits))
	)
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	//获取公钥的数据
	publicKey := privateKey.PublicKey

	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	// X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey) // PKCS1 和 8 是不一致的
	X509PrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return nil, nil, err
	}
	//创建pem.Block结构体对象
	privateBlock := pem.Block{Type: "PRIVATE KEY", Bytes: X509PrivateKey}
	publicBlock := pem.Block{Type: "Public Key", Bytes: X509PublicKey}

	pem.Encode(&privateKeyContentBuffer, &privateBlock)
	pem.Encode(&publicKeyContentBuffer, &publicBlock)

	return privateBlock.Bytes, publicBlock.Bytes, nil
}
