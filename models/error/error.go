package error

import (
	"log"
)

func CheckErr(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
