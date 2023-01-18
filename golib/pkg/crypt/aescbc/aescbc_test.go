package aescbc

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCipherAescbc(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("encode base64 safe", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredEncryptBase64(f.plainText, f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.securedTextBase64, result)
		})
		t.Run("encode hex safe", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredEncryptHex(f.plainText, f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.securedTextHex, result)
		})

		t.Run("decode base64 safe", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredDecryptBase64(f.securedTextBase64, f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.plainText, result)
		})

		t.Run("decode hex safe", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredDecryptHex(f.securedTextHex, f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.plainText, result)
		})
		t.Run("encode base64", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredEncryptBase64(f.plainText, !f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.securedTextBase64NotSafe, result)
		})
		t.Run("encode hex", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredEncryptHex(f.plainText, !f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.securedTextHexNotSafe, result)
		})
		t.Run("decode base64", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredDecryptBase64(f.securedTextBase64NotSafe, !f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.plainText, result)
		})

		t.Run("decode hex safe", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result, err := cipher.SecuredDecryptHex(f.securedTextHexNotSafe, !f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.plainText, result)
		})
		t.Run("get iv", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)
			data, _ := base64.StdEncoding.DecodeString(f.securedTextBase64)

			// act
			result, err := cipher.GetIvBase64FromSecuredText(data)

			// assert
			require.NoError(t, err)
			assert.Equal(t, f.iv, result)
		})

		t.Run("is base64", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result := cipher.IsBase64(f.securedTextBase64)

			// assert
			assert.Equal(t, true, result)
		})

		t.Run("is hex", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, f.iv)

			// act
			result := cipher.IsBase64(f.securedTextHex)

			// assert
			assert.Equal(t, false, result)
		})

		t.Run("set iv", func(t *testing.T) {
			// arrange
			f := setUp(t)

			cipher := New(f.key1, f.key2, "")

			// act
			err := cipher.SetIv(f.iv)

			// assert
			require.NoError(t, err)
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Run("encode base64 iv empty", func(t *testing.T) {
			// arrange
			f := setUp(t)

			plainText := f.plainText
			cipher := New(f.key1, f.key2, "")

			// act
			securedTextBase64, err := cipher.SecuredEncryptBase64(plainText, f.isSafeMode)

			// assert
			require.NoError(t, err)
			assert.NotEqual(t, f.securedTextBase64, securedTextBase64)

		})
		t.Run("encode base64 plaintext empty", func(t *testing.T) {
			// arrange
			f := setUp(t)

			plainText := ""
			cipher := New(f.key1, f.key2, "")

			// act
			securedTextBase64, err := cipher.SecuredEncryptBase64(plainText, f.isSafeMode)

			// assert
			if assert.Error(t, err) {
				assert.EqualError(t, err, "incoming data is empty")
			}
			assert.Equal(t, "", securedTextBase64)

		})
	})
}
