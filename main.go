package main

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
	"os"
	"runtime"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	//InfoLogger.Println("Starting the application...")
	//InfoLogger.Println("Something noteworthy happened")
	//WarningLogger.Println("There is something you should know about")
	//ErrorLogger.Println("Something went wrong")

	// Get Machine ID in var mId
	mId, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	// Get Os in const myOS
	const myOS = runtime.GOOS

	//Unique ID combination of mId (Machine ID) + myOS (IDOS)
	var uniqueId string = myOS + mId

	// Create a key based on OS and machine ID es: windows7fd705af-1b77-42c0-9f00-42330d32e19d
	var length = len([]rune(uniqueId))
	var myOSandMyID string

	// Ensure that the length is 32 byte
	if length < 32 {
		myOSandMyID = addString(uniqueId, length)
	} else {
		myOSandMyID = cutString(uniqueId, 32)
	}

	// Generate a special Key
	data := []byte(myOSandMyID)
	b := md5.Sum(data)
	key := hex.EncodeToString(b[:])

	if os.Args[1] == "ENC" {
		// fmt.Println("CRYPT")
		encrypted := encryptString(os.Args[2], key)
		fmt.Printf("encrypted : %s\n", encrypted)
	} else if os.Args[1] == "DEC" {
		// fmt.Println("DECRYPT")
		decrypted := decryptString(os.Args[2], key)
		fmt.Printf("%s\n", decrypted)
	} else {
		log.Fatal("ONLY ENC AND DEC VALUE ACCEPTED")
	}
}

func encryptString(stringToEncrypt string, keyString string) (encryptedString string) {
	//decode string to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//New Cipher Block created from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//New GCM created
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decryptString(encryptedString string, keyString string) (decryptedString string) {
	//decode string to bytes
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//New Cipher Block created from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//New GCM created
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from enc data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//DENC the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func addString(s string, length int) string {
	var toAdd string
	// Check string
	for i := length; i < 32; i++ {
		toAdd += "$"
	}
	return s + toAdd
}

func cutString(s string, i int) string {
	runes := []rune(s)
	// Check string and cut
	if len(runes) > i {
		return string(runes[:i])
	}
	return s
}
