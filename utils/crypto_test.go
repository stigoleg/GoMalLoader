package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestAESDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")          // 32 bytes for AES-256
	plaintext := []byte("this is a test messagethisisatestbl") // 30 bytes
	paddedPlaintext := pkcs7Pad(plaintext, aes.BlockSize)      // will pad to 32 bytes
	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = byte(i)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("failed to create cipher: %v", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(encrypted, paddedPlaintext)
	ciphertext := append(iv, encrypted...)

	t.Run("correct key", func(t *testing.T) {
		decrypted, err := AESDecrypt(ciphertext, key)
		if err != nil {
			t.Fatalf("AESDecrypt failed: %v", err)
		}
		unpadded := pkcs7Unpad(decrypted)
		if !bytes.Equal(unpadded, plaintext) {
			t.Errorf("decrypted text does not match original. got %x, want %x", unpadded, plaintext)
		}
	})

	t.Run("wrong key", func(t *testing.T) {
		wrongKey := []byte("abcdef0123456789abcdef0123456789")
		_, _ = AESDecrypt(ciphertext, wrongKey) // Should not panic, but output will be garbage
	})
}

func TestAESDecrypt_EdgeCases(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef") // 32 bytes
	iv := make([]byte, 8)                             // Wrong IV size
	plaintext := []byte("test message for edge cases")
	ciphertext := append(iv, plaintext...)
	t.Run("wrong IV size", func(t *testing.T) {
		_, err := AESDecrypt(ciphertext, key)
		if err == nil {
			t.Error("expected error for wrong IV size")
		}
	})

	iv = make([]byte, 16)
	ciphertext = append(iv, []byte{1, 2, 3}...) // Truncated ciphertext
	t.Run("truncated ciphertext", func(t *testing.T) {
		_, err := AESDecrypt(ciphertext, key)
		if err == nil {
			t.Error("expected error for truncated ciphertext")
		}
	})

	iv = make([]byte, 16)
	ciphertext = append(iv, plaintext...)
	shortKey := []byte("shortkey")
	t.Run("invalid key length", func(t *testing.T) {
		_, err := AESDecrypt(ciphertext, shortKey)
		if err == nil {
			t.Error("expected error for invalid key length")
		}
	})
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, pad...)
}

func pkcs7Unpad(data []byte) []byte {
	padLen := int(data[len(data)-1])
	return data[:len(data)-padLen]
}
