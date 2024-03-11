package main

import (
	"encdenc/lib"
	"fmt"
	"log"
	"os"
	"os/user"
	"syscall"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

const pathLog = "/opt/frm/writable/logs/"

func main() {

	file, err := os.OpenFile(pathLog+"enc_dec.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)

	if err != nil {
		fmt.Println("Errore:", err)
		return
	}
	defer file.Close()

	// FileInfo
	fileInfo, err := os.Stat(pathLog + "enc_dec.log")

	// CHECK permission log file
	// GET permission file and check if 664 is correct
	if fmt.Sprintf("%o", fileInfo.Mode().Perm()) != "664" {

		// SET THE LOG WRITE FILE
		log.SetOutput(file)
		log.Println("WARNING ::: Current file has no correct 664 permission. I'll Try to correct it, CONTINUE")

		// GET owner file
		owner, err := user.LookupId(fmt.Sprint(fileInfo.Sys().(*syscall.Stat_t).Uid))
		if err != nil {
			fmt.Println("Error for the owner:", err)
			return
		}
		// fmt.Printf("Owner: %s\n", owner.Username)

		// Ottenere le informazioni sull'utente corrente
		userCurrent, err := user.Current()
		if err != nil {
			fmt.Println("Error for the curretn user:", err)
			return
		}
		// Stampa il nome utente corrente
		// fmt.Printf("Nome Utente Corrente: %s\n", userCurrent.Username)
		if owner.Username == userCurrent.Username {
			log.Println("INFO ::: CurrentUser and Owner file is the same, CONTINUE")
			// Permission change
			err = os.Chmod(pathLog+"enc_dec.log", 0664)

			if err != nil {
				log.Println("ERROR ::: CHANGE PERMISSION LOG FILE", err)
				log.Fatal(err)
			} else {
				log.Println("INFO ::: OK ::: CHANGE PERMISSION LOG FILE COMMITTED NO ERRORS  :: END")
			}
		} else {
			log.Println("WARNING ::: CurrentUser and Owner file is not the same")
		}

	}

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
