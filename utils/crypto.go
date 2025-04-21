package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AESDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize || len(data)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext not full blocks or too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)
	return decrypted, nil
}
