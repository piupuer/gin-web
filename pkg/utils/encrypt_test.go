package utils

import (
	"fmt"
	"testing"
)

func TestRSA(t *testing.T) {
	privateBytes, publicBytes, err := RSAGenKey(4096)
	fmt.Println(string(privateBytes), string(publicBytes), err)
	// 用私钥加密，公钥匙解密
	encodeData, err := RSAEncrypt([]byte("123456"), privateBytes)
	fmt.Println(string(encodeData), err)
	decodeData, err := RSADecrypt(encodeData, publicBytes)
	fmt.Println(string(decodeData), err)

	// 用公钥加密，私钥匙解密
	encodeData2, err := RSAEncrypt([]byte("123456"), publicBytes)
	fmt.Println(string(encodeData2), err)
	decodeData2, err := RSADecrypt(encodeData2, privateBytes)
	fmt.Println(string(decodeData2), err)

}
