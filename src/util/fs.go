package util

import (
	"log"
	"os"
	"path"
)

func GetDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "."
	}
	return path.Join(cwd)
}
