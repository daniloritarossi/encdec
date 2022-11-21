package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"io"
	"log"

	"runtime"
)

func GenerateKey() string {
	// Get Machine ID in var mId
	mId, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	// Get Os in const myOS
	const myOS = runtime.GOOS

	// jh const
	const j = "jY-1"

	// Unique ID combination of mId (Machine ID) + myOS (IDOS)
	var uniqueId string = j + mId + myOS

	// Create a key based on OS and machine ID es: windows7fd705af-1b77-42c0-9f00-42330d32e19d
	var length = len([]rune(uniqueId))
	var myOSandMyID string

	// Ensure that the length is 32 byte
	if length < 32 {
		myOSandMyID = AddString(uniqueId, length)
	} else {
		myOSandMyID = CutString(uniqueId, 32)
	}

	// Generate a special Key
	data := []byte(myOSandMyID)
	b := md5.Sum(data)
	key := hex.EncodeToString(b[:])
	return key
}

func EncryptString(stringToEncrypt string, keyString string) (encryptedString string) {
	// decode string to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	// New Cipher Block created from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// New GCM created
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	// Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func DecryptString(encryptedString string, keyString string) (decryptedString string) {
	// decode string to bytes
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	// New Cipher Block created from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// New GCM created
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from enc data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// DENC the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func AddString(s string, length int) string {
	var toAdd string
	// Check string
	for i := length; i < 32; i++ {
		toAdd += "$"
	}
	return s + toAdd
}

func CutString(s string, i int) string {
	runes := []rune(s)
	// Check string and cut
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}
