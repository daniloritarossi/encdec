package main

import (
	"./lib"
	"fmt"
	"log"
	"os"
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
	key := lib.GenerateKey()
	if os.Args[1] == "ENC" {
		encrypted := lib.EncryptString(os.Args[2], key)
		fmt.Printf("encrypted : %s\n", encrypted)
	} else if os.Args[1] == "DEC" {
		decrypted := lib.DecryptString(os.Args[2], key)
		fmt.Printf("%s\n", decrypted)
	} else {
		log.Fatal("ONLY ENC AND DEC VALUE ACCEPTED")
	}
}
