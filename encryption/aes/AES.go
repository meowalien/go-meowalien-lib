package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func Encrypt(key []byte, plainText []byte) (cipherText []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cipherText = make([]byte, block.BlockSize()+len(plainText))
	iv := cipherText[:block.BlockSize()]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[block.BlockSize():], plainText)

	return
}

func Decrypt(key []byte, cipherText []byte) (decodedmess []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	if len(cipherText) < block.BlockSize() {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[: block.BlockSize()]
	cipherText = cipherText[ block.BlockSize():]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)
	decodedmess = cipherText
	return
}