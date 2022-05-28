package util

import (
	"log"
	"os"
)

func GetDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "."
	}
	return path
}
