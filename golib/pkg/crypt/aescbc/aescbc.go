//go:generate mockgen -source=./aescbc.go -destination=./mocks/aescbc.go -package=mock_aescbc
package aescbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/sha3"
)

type CipherAesCbc struct {
	key1 []byte
	key2 []byte
	iv   []byte
}

type CipherAesCbcInterface interface {
	SecuredEncryptBase64(data string, isSafeMode bool) (string, error)
	SecuredEncryptHex(data string, isSafeMode bool) (string, error)
	SecuredDecryptBase64(data string, isSafeMode bool) (string, error)
	SecuredDecryptHex(data string, isSafeMode bool) (string, error)
	GetIvBase64FromSecuredText(input []byte) (string, error)
	IsBase64(text string) bool
	SetIv(iv string) error
}

func New(key1, key2, iv string) CipherAesCbcInterface {
	firstKey, _ := base64.StdEncoding.DecodeString(key1)
	secondKey, _ := base64.StdEncoding.DecodeString(key2)
	ivBin, _ := base64.StdEncoding.DecodeString(iv)

	return &CipherAesCbc{
		key1: firstKey,
		key2: secondKey,
		iv:   ivBin,
	}
}

func (c *CipherAesCbc) SecuredEncryptBase64(data string, isSafeMode bool) (string, error) {
	out, err := c.securedEncrypt(data, isSafeMode)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func (c *CipherAesCbc) SecuredEncryptHex(data string, isSafeMode bool) (string, error) {
	out, err := c.securedEncrypt(data, isSafeMode)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(out), nil
}

func (c *CipherAesCbc) SecuredDecryptBase64(data string, isSafeMode bool) (string, error) {
	cipherData, _ := base64.StdEncoding.DecodeString(data)
	out, err := c.securedDecrypt(cipherData, isSafeMode)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *CipherAesCbc) SecuredDecryptHex(data string, isSafeMode bool) (string, error) {
	cipherData, _ := hex.DecodeString(data)
	out, err := c.securedDecrypt(cipherData, isSafeMode)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Encrypts data according to AES-CBC algorithm using key key1.
// Generates hmac hash of aes-cbc encrypted data based on sha3-512 algorithm using
// key2. If key2 is empty then do not generate sha3-512 hash.
// If iv is provided then uses it to crypt data otherwise the function
// will generate random iv (recommended).
// Returns concatenation of iv, sha3 hash of aes-cbc encoded data (if key2 is presented), and aes-cbc encoded data.
func (c *CipherAesCbc) securedEncrypt(data string, isSafeMode bool) ([]byte, error) {
	var firstEncrypted, secondEncrypted []byte

	if data == "" {
		return nil, errors.New("incoming data is empty")
	}

	plaintext := []byte(data)

	block, err := aes.NewCipher(c.key1)
	if err != nil {
		return nil, err
	}

	plaintext = pkcs7Padding(plaintext, block.BlockSize())
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, errors.New("plaintext is not a multiple of the block size")
	}

	firstEncrypted = make([]byte, aes.BlockSize+len(plaintext))
	if len(c.iv) == 0 {
		c.iv = firstEncrypted[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, c.iv); err != nil {
			return nil, err
		}
	} else {
		copy(firstEncrypted, c.iv)
	}

	mode := cipher.NewCBCEncrypter(block, c.iv)
	mode.CryptBlocks(firstEncrypted[aes.BlockSize:], plaintext)

	if isSafeMode {
		secondEncrypted = hmacSha3(firstEncrypted[aes.BlockSize:], c.key2)
	}

	out := append(c.iv, secondEncrypted...)
	out = append(out, firstEncrypted[aes.BlockSize:]...)

	return out, nil
}

// Splits incoming data by 3 parts:
// - iv,
// - second encrypted data (sha3-512 hmac hash of aes-cbc encrypted data) if key is presented,
// - and first encrypted data (aes-cbc encrypted data).
// Decrypts first encrypted data to plain text using iv and key1.
// If key2 is present then perform timing attack safe comparison:
// - create second encrypted new of first encrypted data
// - compares with second encrypted data
// Returns encrypted plain text if success, otherwise empty string and corresponding error.
func (c *CipherAesCbc) securedDecrypt(input []byte, isSafeMode bool) (string, error) {
	var secondSize int
	var secondEncrypted []byte

	if len(input) == 0 {
		return "", errors.New("incoming data is empty")
	}

	iv := make([]byte, aes.BlockSize)
	if isSafeMode {
		secondSize = 64
		secondEncrypted = make([]byte, len(input[aes.BlockSize:aes.BlockSize+secondSize]))
	}

	firstEncrypted := make([]byte, len(input[aes.BlockSize+secondSize:]))
	copy(iv, input[:aes.BlockSize])
	if isSafeMode {
		copy(secondEncrypted, input[aes.BlockSize:aes.BlockSize+secondSize])
	}
	copy(firstEncrypted, input[aes.BlockSize+secondSize:])

	block, err := aes.NewCipher(c.key1)
	if err != nil {
		return "", err
	}

	if len(firstEncrypted) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	if len(firstEncrypted)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plainText := make([]byte, len(firstEncrypted)+len(iv))

	in := append(iv, firstEncrypted...)
	mode.CryptBlocks(plainText, in)

	if isSafeMode {
		secondEncryptedNew := hmacSha3(firstEncrypted, c.key2)
		if !bytes.Equal(secondEncrypted, secondEncryptedNew) {
			return "", errors.New("second encrypted not equal to second encrypted new")
		}
	}
	plainText = pkcs7UnPadding(plainText[aes.BlockSize:])

	return string(plainText), nil
}

func (c *CipherAesCbc) GetIvBase64FromSecuredText(input []byte) (string, error) {
	if len(input) == 0 {
		return "", errors.New("incoming data is empty")
	}

	iv := make([]byte, aes.BlockSize)
	copy(iv, input[:aes.BlockSize])

	return base64.StdEncoding.EncodeToString(iv), nil

}

func (c *CipherAesCbc) IsBase64(text string) bool {
	var isBase64 bool
	dst := make([]byte, hex.DecodedLen(len(text)))
	if _, err := hex.Decode(dst, []byte(text)); err != nil {
		isBase64 = true
	}
	return isBase64
}

func (c *CipherAesCbc) SetIv(iv string) error {
	var err error
	c.iv, err = base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return err
	}
	return nil
}

func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func hmacSha3(text, key []byte) []byte {
	mac := hmac.New(sha3.New512, key)
	mac.Write(text)
	encrypted := mac.Sum(nil)
	return encrypted
}
